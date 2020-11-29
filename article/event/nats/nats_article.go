package nats

import (
	"bytes"
	"context"
	"encoding/gob"

	"github.com/arifseft/article-api/domain"
	"github.com/nats-io/nats.go"
)

const (
	topic = "article:created"
)

type natsArticleEvent struct {
	Conn         *nats.Conn
	Subscription *nats.Subscription
}

func NewNatsArticleEvent(Conn *nats.Conn) domain.ArticleEvent {
	return &natsArticleEvent{Conn: Conn}
}

func (n *natsArticleEvent) PublishArticleCreated(ctx context.Context, a domain.Article) (err error) {
	b := bytes.Buffer{}
	err = gob.NewEncoder(&b).Encode(a)
	if err != nil {
		return err
	}

	return n.Conn.Publish(topic, b.Bytes())
}

func (n *natsArticleEvent) SubscribeArticleCreated(ctx context.Context, f func(domain.Article)) (err error) {
	a := domain.Article{}
	n.Subscription, err = n.Conn.Subscribe(topic, func(msg *nats.Msg) {
		b := bytes.Buffer{}
		b.Write(msg.Data)
		gob.NewDecoder(&b).Decode(&a)
		f(a)
	})
	return
}
