/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-06-05 16:29:54
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-16 11:48:55
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package flake_test

import (
	"sync"
	"testing"
	"time"

	"github.com/Mrzhouyl/go-aisdk/flake"
)

func TestNextRequestID(t *testing.T) {
	flake, err := flake.New(flake.Settings{})
	if err != nil {
		t.Fatalf("new flake error: %v", err)
	}

	requestId, err := flake.RequestID()
	if err != nil {
		t.Fatalf("next request id error: %v", err)
	}
	t.Logf("next request id: %v", requestId)
}

func TestNextRequestIDUniqueness(t *testing.T) {
	flake, err := flake.New(flake.Settings{})
	if err != nil {
		t.Fatalf("new flake error: %v", err)
	}

	const numIDs = 1000000
	ids := make(map[string]bool)
	var wg sync.WaitGroup
	var mutex sync.Mutex
	duplicateCount := 0
	start := time.Now()

	for range numIDs {
		wg.Add(1)
		go func() {
			defer wg.Done()

			id, err := flake.RequestID()
			if err != nil {
				t.Errorf("next request id error: %v", err)
				return
			}

			mutex.Lock()
			if ids[id] {
				duplicateCount++
				t.Errorf("duplicate id: %v", id)
			} else {
				ids[id] = true
			}
			mutex.Unlock()
		}()
	}

	wg.Wait()
	elapsed := time.Since(start)

	t.Logf("generate count: %d", numIDs)
	t.Logf("unique count: %d", len(ids))
	t.Logf("duplicate count: %d", duplicateCount)
	t.Logf("elapsed: %v", elapsed)
	t.Logf("QPS: %.0f", float64(numIDs)/elapsed.Seconds())

	if duplicateCount == 0 {
		t.Logf("✅ all ids are unique, test passed!")
	}
}
