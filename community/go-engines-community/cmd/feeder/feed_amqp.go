package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

func (f *Feeder) sendIterable(ctx context.Context, iterable []interface{}) error {
	for i, v := range iterable {
		f.logger.Info().Msgf("sending event %d/%d", i+1, len(iterable))
		_, ok := v.(map[string]interface{})

		if !ok {
			return fmt.Errorf("sending event: not a map[string]interface{}")
		}

		bv, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("sending event: %v", err)
		}
		if err = f.sendBytes(ctx, bv, "#"); err != nil {
			return fmt.Errorf("sending event: %v", err)
		}
	}

	return nil
}

func (f *Feeder) sendLoop(ctx context.Context, content []byte) error {
	var err error
	if f.flags.CheckJSON {
		var ref interface{}
		err = json.Unmarshal(content, &ref)

		if err != nil {
			return fmt.Errorf("cannot unmarshal: %v", err)
		}

		var iterable bool
		switch ref := ref.(type) {
		case map[string]interface{}:
			f.logger.Info().Msgf("sending one event from file %s", f.flags.File)

		case []interface{}:
			length := len(ref)
			f.logger.Info().Msgf("sending %d events from file %s", length, f.flags.File)
			iterable = true
			err = f.sendIterable(ctx, ref)
		}

		if !iterable {
			err = f.sendBytes(ctx, content, "#")
		}
	} else {
		err = f.sendBytes(ctx, content, "#")
	}

	return err
}

func (f *Feeder) modeSendEvent(ctx context.Context) error {
	if err := f.setupAmqp(); err != nil {
		return err
	}

	fi, err := os.Stat(f.flags.File)
	if err != nil {
		return err
	}

	if fi.IsDir() {
		return fmt.Errorf("file %s is a directory", f.flags.File)
	}

	content, err := os.ReadFile(f.flags.File)

	if err != nil {
		return fmt.Errorf("reading file %s: %v", f.flags.File, err)
	}

	sendLoop := true

	for sendLoop {
		err := f.sendLoop(ctx, content)

		if err != nil {
			return fmt.Errorf("sending event: %v", err)
		}

		// breaking the loop
		sendLoop = sendLoop && f.flags.Loop
	}

	f.logger.Info().Msg("event sent.")

	return err
}
