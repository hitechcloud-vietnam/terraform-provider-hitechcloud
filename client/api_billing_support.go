// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
)

// Billing & Contracts (folder "Billing & Contracts") extras.

// ApplyInvoiceCredit applies account credit to an invoice
// (POST /api/invoice/:id/credit).
func (c *Client) ApplyInvoiceCredit(ctx context.Context, invoiceID, amount string) (map[string]any, error) {
	return c.postMap(ctx, "/api/invoice/"+PathEscape(invoiceID)+"/credit", Query("amount", amount))
}

// Support extras: ticket attachments, news and knowledgebase (folder
// "Support").

// GetTicketAttachment downloads a ticket attachment
// (GET /api/ticket/attachment/:file?number=).
func (c *Client) GetTicketAttachment(ctx context.Context, number, file string) ([]byte, error) {
	return c.GetRaw(ctx, "/api/ticket/attachment/"+PathEscape(file), Query("number", number))
}

// ListNews lists portal news (GET /api/news).
func (c *Client) ListNews(ctx context.Context) ([]any, error) {
	return c.getList(ctx, "/api/news", nil)
}

// GetNewsItem returns a single news item (GET /api/news/:news_id).
func (c *Client) GetNewsItem(ctx context.Context, newsID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/news/"+PathEscape(newsID), nil)
}

// ListKnowledgebaseCategories lists knowledgebase categories
// (GET /api/knowledgebase).
func (c *Client) ListKnowledgebaseCategories(ctx context.Context) ([]any, error) {
	return c.getList(ctx, "/api/knowledgebase", nil)
}

// GetKnowledgebaseCategory returns a knowledgebase category
// (GET /api/knowledgebase/:category_id).
func (c *Client) GetKnowledgebaseCategory(ctx context.Context, categoryID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/knowledgebase/"+PathEscape(categoryID), nil)
}

// GetKnowledgebaseArticle returns a knowledgebase article
// (GET /api/knowledgebase/article/:article_id).
func (c *Client) GetKnowledgebaseArticle(ctx context.Context, articleID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/knowledgebase/article/"+PathEscape(articleID), nil)
}

// Notifications extras (folder "Notifications").

// ListNewNotifications lists unacknowledged portal notifications
// (GET /api/notifications/new).
func (c *Client) ListNewNotifications(ctx context.Context, relType, relID string) ([]any, error) {
	q := Query()
	if relType != "" {
		q.Set("rel_type", relType)
	}
	if relID != "" {
		q.Set("rel_id", relID)
	}
	return c.getList(ctx, "/api/notifications/new", q)
}

// AcknowledgeNotification marks a notification as read
// (PUT /api/notifications/:id/ack).
func (c *Client) AcknowledgeNotification(ctx context.Context, id string) error {
	_, err := c.putMap(ctx, "/api/notifications/"+PathEscape(id)+"/ack", nil)
	return err
}

// GetStatusDetails returns the details of a service status entry
// (PUT /api/statuses/:id — the API exposes this lookup as PUT).
func (c *Client) GetStatusDetails(ctx context.Context, id string) (map[string]any, error) {
	return c.putMap(ctx, "/api/statuses/"+PathEscape(id), nil)
}
