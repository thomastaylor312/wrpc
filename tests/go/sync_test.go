//go:generate $WIT_BINDGEN_WRPC go --gofmt=false --world sync-client --out-dir bindings/sync_client --package wrpc.io/tests/go/bindings/sync_client ../wit

package integration_test

import (
	"context"
	"log/slog"
	"reflect"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	wrpc "wrpc.io/go"
	wrpcnats "wrpc.io/go/nats"
	integration "wrpc.io/tests/go"
	"wrpc.io/tests/go/bindings/sync_client/foo"
	"wrpc.io/tests/go/bindings/sync_client/wrpc_test/integration/sync"
	"wrpc.io/tests/go/bindings/sync_server"
	"wrpc.io/tests/go/internal"
)

func TestSync(t *testing.T) {
	natsSrv := internal.RunNats(t)
	nc, err := nats.Connect(natsSrv.ClientURL())
	if err != nil {
		t.Errorf("failed to connect to NATS.io: %s", err)
		return
	}
	defer nc.Close()
	defer func() {
		if err := nc.Drain(); err != nil {
			t.Errorf("failed to drain NATS.io connection: %s", err)
			return
		}
	}()
	client := wrpcnats.NewClient(nc, wrpcnats.WithPrefix("go"))

	var h integration.SyncHandler
	stop, err := sync_server.Serve(client, h, h)
	if err != nil {
		t.Errorf("failed to serve `sync-server` world: %s", err)
		return
	}

	var cancel func()
	ctx := context.Background()
	dl, ok := t.Deadline()
	if ok {
		ctx, cancel = context.WithDeadline(ctx, dl)
	} else {
		ctx, cancel = context.WithTimeout(ctx, time.Minute)
	}
	defer cancel()

	{
		slog.DebugContext(ctx, "calling `wrpc-test:integration/sync-client.foo.f`")
		v, err := foo.F(ctx, client, "f")
		if err != nil {
			t.Errorf("failed to call `wrpc-test:integration/sync-client.foo.f`: %s", err)
			return
		}
		if v != 42 {
			t.Errorf("expected: 42, got: %d", v)
			return
		}
	}
	{
		slog.DebugContext(ctx, "calling `wrpc-test:integration/sync-client.foo.foo`")
		err := foo.Foo(ctx, client, "foo")
		if err != nil {
			t.Errorf("failed to call `wrpc-test:integration/sync-client.foo.foo`: %s", err)
			return
		}
	}
	{
		slog.DebugContext(ctx, "calling `wrpc-test:integration/sync.fallible`")
		v, err := sync.Fallible(ctx, client, true)
		if err != nil {
			t.Errorf("failed to call `wrpc-test:integration/sync.fallible`: %s", err)
			return
		}
		expected := wrpc.Ok[string](true)
		if !reflect.DeepEqual(v, expected) {
			t.Errorf("expected: %#v, got: %#v", expected, v)
			return
		}
	}
	{
		slog.DebugContext(ctx, "calling `wrpc-test:integration/sync.fallible`")
		v, err := sync.Fallible(ctx, client, false)
		if err != nil {
			t.Errorf("failed to call `wrpc-test:integration/sync.fallible`: %s", err)
			return
		}
		expected := wrpc.Err[bool]("test")
		if !reflect.DeepEqual(v, expected) {
			t.Errorf("expected: %#v, got: %#v", expected, v)
			return
		}
	}
	{
		slog.DebugContext(ctx, "calling `wrpc-test:integration/sync.numbers`")
		v0, v1, v2, v3, v4, v5, v6, v7, v8, v9, err := sync.Numbers(ctx, client)
		if err != nil {
			t.Errorf("failed to call `wrpc-test:integration/sync.numbers`: %s", err)
			return
		}
		if v0 != 1 {
			t.Errorf("expected: 1, got: %#v", v0)
			return
		}
		if v1 != 2 {
			t.Errorf("expected: 2, got: %#v", v1)
			return
		}
		if v2 != 3 {
			t.Errorf("expected: 3, got: %#v", v2)
			return
		}
		if v3 != 4 {
			t.Errorf("expected: 4, got: %#v", v3)
			return
		}
		if v4 != 5 {
			t.Errorf("expected: 5, got: %#v", v4)
			return
		}
		if v5 != 6 {
			t.Errorf("expected: 6, got: %#v", v5)
			return
		}
		if v6 != 7 {
			t.Errorf("expected: 7, got: %#v", v6)
			return
		}
		if v7 != 8 {
			t.Errorf("expected: 8, got: %#v", v7)
			return
		}
		if v8 != 9 {
			t.Errorf("expected: 9, got: %#v", v8)
			return
		}
		if v9 != 10 {
			t.Errorf("expected: 10, got: %#v", v9)
			return
		}
	}
	{
		slog.DebugContext(ctx, "calling `wrpc-test:integration/sync.with-flags`")
		v, err := sync.WithFlags(ctx, client, true, false, true)
		if err != nil {
			t.Errorf("failed to call `wrpc-test:integration/sync.with-flags`: %s", err)
			return
		}
		expected := &sync.Abc{A: true, B: false, C: true}
		if !reflect.DeepEqual(v, expected) {
			t.Errorf("expected: %v, got: %#v", expected, v)
			return
		}
	}
	{
		v, err := sync.WithVariantOption(ctx, client, true)
		if err != nil {
			t.Errorf("failed to call `wrpc-test:integration/sync.with-variant-option`: %s", err)
			return
		}
		expected := sync.NewVarVar(&sync.Rec{
			Nested: &sync.RecNested{
				Foo: "bar",
			},
		})
		if !reflect.DeepEqual(v, expected) {
			t.Errorf("expected: %v, got: %#v", expected, v)
			return
		}
	}
	{
		v, err := sync.WithVariantOption(ctx, client, false)
		if err != nil {
			t.Errorf("failed to call `wrpc-test:integration/sync.with-variant-option`: %s", err)
			return
		}
		var expected *sync.Var
		if !reflect.DeepEqual(v, expected) {
			t.Errorf("expected: %v, got: %#v", expected, v)
			return
		}
	}
	{
		v, err := sync.WithVariantList(ctx, client)
		if err != nil {
			t.Errorf("failed to call `wrpc-test:integration/sync.with-variant-list`: %s", err)
			return
		}
		expected := []*sync.Var{
			sync.NewVarEmpty(),
			sync.NewVarVar(&sync.Rec{
				Nested: &sync.RecNested{
					Foo: "foo",
				},
			}),
			sync.NewVarEmpty(),
			sync.NewVarEmpty(),
			sync.NewVarEmpty(),
			sync.NewVarVar(&sync.Rec{
				Nested: &sync.RecNested{
					Foo: "bar",
				},
			}),
		}
		if !reflect.DeepEqual(v, expected) {
			t.Errorf("expected: %v, got: %#v", expected, v)
			return
		}
	}
	{
		v, err := sync.WithRecord(ctx, client)
		if err != nil {
			t.Errorf("failed to call `wrpc-test:integration/sync.with-record`: %s", err)
			return
		}
		expected := &sync.Rec{
			Nested: &sync.RecNested{
				Foo: "foo",
			},
		}
		if !reflect.DeepEqual(v, expected) {
			t.Errorf("expected: %v, got: %#v", expected, v)
			return
		}
	}
	{
		v, err := sync.WithRecordList(ctx, client, 3)
		if err != nil {
			t.Errorf("failed to call `wrpc-test:integration/sync.with-record-list`: %s", err)
			return
		}
		expected := []*sync.Rec{
			{
				Nested: &sync.RecNested{
					Foo: "0",
				},
			},
			{
				Nested: &sync.RecNested{
					Foo: "1",
				},
			},
			{
				Nested: &sync.RecNested{
					Foo: "2",
				},
			},
		}
		if !reflect.DeepEqual(v, expected) {
			t.Errorf("expected: %v, got: %#v", expected, v)
			return
		}
	}

	if err = stop(); err != nil {
		t.Errorf("failed to stop serving `sync-server` world: %s", err)
		return
	}
	if nc.NumSubscriptions() != 0 {
		t.Errorf("NATS subscriptions leaked: %d active after client stop", nc.NumSubscriptions())
	}
}
