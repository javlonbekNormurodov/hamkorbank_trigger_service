package trigger_listener_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hamkorbank/config"
	"hamkorbank/pkg/logger"
	"syscall"
)

func (t *triggerListener) Listen(ctx context.Context, data []byte) error {
	var resp = &Message{}
	err := json.Unmarshal(data, resp)
	if err != nil {
		t.log.Error("error while consuming ", logger.Error(err))
		return err
	}

	t.log.Info("Debug", logger.Any("resp ", resp))

	err = t.rabbitmq.Publish(ctx, config.AllDebug, data)
	if err != nil {
		t.log.Error("Error while publishing data", logger.Error(err))
	}

	path := fmt.Sprintf("http://%s:%d/v1/phone/%s", t.cfg.RestServiceHost, t.cfg.RestServicePort, resp.RecordId)

	body, status, err := t.httpClient.Request("GET", path, "application/json", "", nil, "")
	if errors.Is(err, syscall.ECONNREFUSED) {
		panic(err)
	}

	var b []byte

	if status == 404 {
		b, err = json.Marshal(NotFound{
			NotFound: resp.RecordId,
		})
		if err != nil {
			t.log.Error("Error while marshaling data", logger.Error(err))
		}

		err = t.rabbitmq.Publish(ctx, config.AllErrors, b)
		if err != nil {
			t.log.Error("Error while publishing data", logger.Error(err))
			return err
		}
		return nil
	}

	if status == 500 {
		b, err = json.Marshal(Message{
			RecordId: resp.RecordId,
		})
		if err != nil {
			t.log.Error("Error while marshaling data", logger.Error(err))
		}

		err = t.rabbitmq.Publish(ctx, config.Consumer, b)
		if err != nil {
			t.log.Error("Error while publishing data", logger.Error(err))
			return err
		}
		return nil
	}

	var respBody Response
	_ = json.Unmarshal(body, &respBody)
	b, err = json.Marshal(respBody.Data)
	err = t.rabbitmq.Publish(ctx, config.AllInfo, b)
	if err != nil {
		t.log.Error("Error while publishing data", logger.Error(err))
		return err
	}

	return nil
}
