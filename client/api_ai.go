// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// All endpoints in this file belong to the "HiTechCloud AI Factory" folder and
// are scoped by the HostBill service id.

// ---------------------------------------------------------------------------
// GPU instances
// ---------------------------------------------------------------------------

// AIInstance describes a HiTechCloud AI Factory GPU instance.
type AIInstance struct {
	ID                  string
	Name                string
	Status              string
	Cloud               string
	Region              string
	InstanceType        string
	ShadeCloud          bool
	OS                  string
	TemplateID          string
	SSHKeyID            string
	VolumeIDs           []string
	Tags                []string
	Envs                map[string]string
	AutoDelete          bool
	Alert               bool
	LaunchConfiguration string // JSON document
	VolumeMount         string // JSON document
	IP                  string
}

// AIInstanceCreate carries the parameters of POST /api/service/:id/instances.
type AIInstanceCreate struct {
	Cloud               string
	Region              string
	ShadeInstanceType   string
	ShadeCloud          *bool
	Name                string
	OS                  string
	TemplateID          string
	SSHKeyID            string
	VolumeIDs           []string
	LaunchConfiguration string // JSON document
	AutoDelete          *bool
	Alert               *bool
	VolumeMount         string // JSON document
	Tags                []string
	Envs                map[string]string
}

// ListAIInstances lists GPU instances (GET /api/service/:id/instances).
func (c *Client) ListAIInstances(ctx context.Context, serviceID string) ([]AIInstance, error) {
	var raw any
	path := "/api/service/" + PathEscape(serviceID) + "/instances"
	if err := c.Get(ctx, path, nil, &raw); err != nil {
		return nil, err
	}
	list := ExtractList(raw, "instances")
	out := make([]AIInstance, 0, len(list))
	for _, item := range list {
		if m, ok := AsMap(item); ok {
			out = append(out, parseAIInstance(m))
		}
	}
	return out, nil
}

// GetAIInstance fetches one instance (GET /api/service/:id/instances/:instance_id).
func (c *Client) GetAIInstance(ctx context.Context, serviceID, instanceID string) (*AIInstance, error) {
	var raw any
	path := "/api/service/" + PathEscape(serviceID) + "/instances/" + PathEscape(instanceID)
	if err := c.Get(ctx, path, nil, &raw); err != nil {
		return nil, err
	}
	m, ok := AsMap(raw)
	if !ok {
		if list := ExtractList(raw, "instances"); len(list) > 0 {
			m, _ = AsMap(list[0])
		}
	}
	if m == nil {
		return nil, fmt.Errorf("unexpected response shape for AI instance %s", instanceID)
	}
	inst := parseAIInstance(m)
	if inst.ID == "" {
		inst.ID = instanceID
	}
	return &inst, nil
}

// CreateAIInstance provisions a GPU instance
// (POST /api/service/:id/instances) and resolves its id.
func (c *Client) CreateAIInstance(ctx context.Context, serviceID string, in AIInstanceCreate) (string, error) {
	var raw any
	path := "/api/service/" + PathEscape(serviceID) + "/instances"
	q := Query(
		"cloud", in.Cloud,
		"region", in.Region,
		"shade_instance_type", in.ShadeInstanceType,
		"name", in.Name,
		"os", in.OS,
		"template_id", in.TemplateID,
		"ssh_key_id", in.SSHKeyID,
	)
	QueryBool(q, "shade_cloud", in.ShadeCloud)
	QueryBool(q, "auto_delete", in.AutoDelete)
	QueryBool(q, "alert", in.Alert)
	if err := JSONQuery(q, "launch_configuration", json.RawMessage(in.LaunchConfiguration)); err != nil {
		return "", err
	}
	if err := JSONQuery(q, "volume_mount", json.RawMessage(in.VolumeMount)); err != nil {
		return "", err
	}
	if err := JSONQuery(q, "volume_ids", in.VolumeIDs); err != nil {
		return "", err
	}
	if err := JSONQuery(q, "tags", in.Tags); err != nil {
		return "", err
	}
	if err := JSONQuery(q, "envs", in.Envs); err != nil {
		return "", err
	}
	if err := c.Post(ctx, path, q, &raw); err != nil {
		return "", err
	}
	if id := FirstString(raw, "id", "instance_id", "instanceid"); id != "" {
		return id, nil
	}
	// Fallback: locate by name.
	instances, err := c.ListAIInstances(ctx, serviceID)
	if err != nil {
		return "", fmt.Errorf("AI instance was created but its id could not be resolved: %w", err)
	}
	for _, inst := range instances {
		if in.Name != "" && inst.Name == in.Name {
			return inst.ID, nil
		}
	}
	return "", fmt.Errorf("AI instance %q was created but could not be found", in.Name)
}

// UpdateAIInstance updates mutable fields
// (POST /api/service/:id/instances/:instance_id/update).
func (c *Client) UpdateAIInstance(ctx context.Context, serviceID, instanceID, name string, autoDelete, alert *bool, tags []string) error {
	path := "/api/service/" + PathEscape(serviceID) + "/instances/" + PathEscape(instanceID) + "/update"
	q := Query("name", name)
	QueryBool(q, "auto_delete", autoDelete)
	QueryBool(q, "alert", alert)
	if err := JSONQuery(q, "tags", tags); err != nil {
		return err
	}
	return c.Post(ctx, path, q, nil)
}

// DeleteAIInstance deletes an instance
// (POST /api/service/:id/instances/:instance_id/delete).
func (c *Client) DeleteAIInstance(ctx context.Context, serviceID, instanceID string) error {
	path := "/api/service/" + PathEscape(serviceID) + "/instances/" + PathEscape(instanceID) + "/delete"
	return c.Post(ctx, path, nil, nil)
}

// WaitForAIInstanceStable polls the instance until it leaves transient
// provisioning states (or the timeout elapses). The API creates GPU instances
// asynchronously, so Terraform must wait for a terminal state.
func (c *Client) WaitForAIInstanceStable(ctx context.Context, serviceID, instanceID string, timeout time.Duration) (*AIInstance, error) {
	deadline := time.Now().Add(timeout)
	for {
		inst, err := c.GetAIInstance(ctx, serviceID, instanceID)
		if err == nil {
			if aiInstanceReady(inst.Status) {
				return inst, nil
			}
			if aiInstanceFailed(inst.Status) {
				return inst, fmt.Errorf("AI instance %s entered failure state %q", instanceID, inst.Status)
			}
		}
		if time.Now().After(deadline) {
			if err != nil {
				return nil, fmt.Errorf("timed out waiting for AI instance %s: %w", instanceID, err)
			}
			return inst, fmt.Errorf("timed out waiting for AI instance %s to leave state %q", instanceID, inst.Status)
		}
		if err := sleepCtx(ctx, 10*time.Second); err != nil {
			return nil, err
		}
	}
}

// WaitForAIInstanceDeleted polls until the instance is gone (404).
func (c *Client) WaitForAIInstanceDeleted(ctx context.Context, serviceID, instanceID string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		_, err := c.GetAIInstance(ctx, serviceID, instanceID)
		if IsNotFound(err) {
			return nil
		}
		if err == nil && time.Now().After(deadline) {
			return fmt.Errorf("timed out waiting for AI instance %s deletion", instanceID)
		}
		if err != nil && !IsNotFound(err) && time.Now().After(deadline) {
			return fmt.Errorf("timed out waiting for AI instance %s deletion: %w", instanceID, err)
		}
		if err := sleepCtx(ctx, 10*time.Second); err != nil {
			return err
		}
	}
}

var aiTransientStates = map[string]bool{
	"creating": true, "create_in_progress": true, "pending": true,
	"provisioning": true, "initializing": true, "starting": true,
	"updating": true, "deleting": true, "stopping": true, "restarting": true,
	"queued": true, "in_progress": true, "launching": true, "": true,
}

var aiFailedStates = map[string]bool{
	"error": true, "failed": true, "failure": true, "terminated": true,
	"destroyed": true, "cancelled": true, "canceled": true,
}

func aiInstanceReady(status string) bool {
	s := strings.ToLower(strings.TrimSpace(status))
	return !aiTransientStates[s] && !aiFailedStates[s]
}

func aiInstanceFailed(status string) bool {
	return aiFailedStates[strings.ToLower(strings.TrimSpace(status))]
}

func parseAIInstance(m map[string]any) AIInstance {
	inst := AIInstance{
		ID:           FirstString(m, "id", "instance_id", "instanceid"),
		Name:         FirstString(m, "name", "label", "instance_name"),
		Status:       FirstString(m, "status", "state"),
		Cloud:        FirstString(m, "cloud", "provider", "cloud_provider"),
		Region:       FirstString(m, "region", "location"),
		InstanceType: FirstString(m, "shade_instance_type", "instance_type", "type", "machine_type"),
		ShadeCloud:   FirstBool(m, "shade_cloud", "shadecloud"),
		OS:           FirstString(m, "os", "image", "os_image"),
		TemplateID:   asString(First(m, "template_id", "template")),
		SSHKeyID:     asString(First(m, "ssh_key_id", "sshkey_id", "ssh_key")),
		VolumeIDs:    StringList(First(m, "volume_ids", "volumes", "volume_id")),
		Tags:         StringList(First(m, "tags", "tag")),
		Envs:         StringMap(First(m, "envs", "env", "environment", "environment_variables")),
		AutoDelete:   FirstBool(m, "auto_delete", "autodelete"),
		Alert:        FirstBool(m, "alert", "alerts"),
		IP:           FirstString(m, "ip", "public_ip", "address", "ip_address"),
	}
	inst.LaunchConfiguration = jsonFieldString(m, "launch_configuration", "launchConfiguration")
	inst.VolumeMount = jsonFieldString(m, "volume_mount", "volumeMount", "mount")
	return inst
}

// jsonFieldString renders a structured response field back into a canonical
// JSON string so it can be stored in Terraform state verbatim.
func jsonFieldString(m map[string]any, keys ...string) string {
	v := First(m, keys...)
	if v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	default:
		raw, err := json.Marshal(t)
		if err != nil {
			return ""
		}
		return string(raw)
	}
}

// ---------------------------------------------------------------------------
// SSH keys
// ---------------------------------------------------------------------------

// AISSHKey is an OpenSSH key registered with the AI Factory service.
type AISSHKey struct {
	ID        string
	Name      string
	PublicKey string
	IsDefault bool
}

// ListAISSHKeys lists SSH keys (GET /api/service/:id/sshkeys).
func (c *Client) ListAISSHKeys(ctx context.Context, serviceID string) ([]AISSHKey, error) {
	var raw any
	path := "/api/service/" + PathEscape(serviceID) + "/sshkeys"
	if err := c.Get(ctx, path, nil, &raw); err != nil {
		return nil, err
	}
	list := ExtractList(raw, "sshkeys", "ssh_keys", "keys")
	out := make([]AISSHKey, 0, len(list))
	for _, item := range list {
		if m, ok := AsMap(item); ok {
			out = append(out, parseAISSHKey(m))
		}
	}
	return out, nil
}

// GetAISSHKey fetches one key (GET /api/service/:id/sshkeys/:key_id).
func (c *Client) GetAISSHKey(ctx context.Context, serviceID, keyID string) (*AISSHKey, error) {
	var raw any
	path := "/api/service/" + PathEscape(serviceID) + "/sshkeys/" + PathEscape(keyID)
	if err := c.Get(ctx, path, nil, &raw); err != nil {
		return nil, err
	}
	m, ok := AsMap(raw)
	if !ok {
		if list := ExtractList(raw, "sshkeys"); len(list) > 0 {
			m, _ = AsMap(list[0])
		}
	}
	if m == nil {
		return nil, fmt.Errorf("unexpected response shape for SSH key %s", keyID)
	}
	key := parseAISSHKey(m)
	if key.ID == "" {
		key.ID = keyID
	}
	return &key, nil
}

// CreateAISSHKey adds an OpenSSH public key (POST /api/service/:id/sshkeys).
// The API documents both `public_key` and `key` parameter names; both are sent
// with the same value for compatibility.
func (c *Client) CreateAISSHKey(ctx context.Context, serviceID, name, publicKey string) (string, error) {
	var raw any
	path := "/api/service/" + PathEscape(serviceID) + "/sshkeys"
	q := Query("public_key", publicKey, "key", publicKey, "name", name)
	if err := c.Post(ctx, path, q, &raw); err != nil {
		return "", err
	}
	if id := FirstString(raw, "id", "key_id", "keyid"); id != "" {
		return id, nil
	}
	// Fallback: locate by public key.
	keys, err := c.ListAISSHKeys(ctx, serviceID)
	if err != nil {
		return "", fmt.Errorf("SSH key was created but its id could not be resolved: %w", err)
	}
	for _, k := range keys {
		if k.PublicKey == publicKey || (name != "" && k.Name == name) {
			return k.ID, nil
		}
	}
	return "", fmt.Errorf("SSH key %q was created but could not be found", name)
}

// DeleteAISSHKey removes a key (POST /api/service/:id/sshkeys/:key_id/delete).
func (c *Client) DeleteAISSHKey(ctx context.Context, serviceID, keyID string) error {
	path := "/api/service/" + PathEscape(serviceID) + "/sshkeys/" + PathEscape(keyID) + "/delete"
	return c.Post(ctx, path, nil, nil)
}

// SetDefaultAISSHKey marks a key as default
// (POST /api/service/:id/sshkeys/:key_id/setdefault).
func (c *Client) SetDefaultAISSHKey(ctx context.Context, serviceID, keyID string) error {
	path := "/api/service/" + PathEscape(serviceID) + "/sshkeys/" + PathEscape(keyID) + "/setdefault"
	return c.Post(ctx, path, nil, nil)
}

func parseAISSHKey(m map[string]any) AISSHKey {
	return AISSHKey{
		ID:        FirstString(m, "id", "key_id", "keyid"),
		Name:      FirstString(m, "name", "label", "key_name"),
		PublicKey: FirstString(m, "public_key", "key", "publickey", "public_key_data"),
		IsDefault: FirstBool(m, "is_default", "default", "isdefault"),
	}
}

// ---------------------------------------------------------------------------
// Volumes
// ---------------------------------------------------------------------------

// AIVolume is a storage volume.
type AIVolume struct {
	ID       string
	Name     string
	Cloud    string
	Region   string
	SizeInGB int64
	Status   string
}

// ListAIVolumes lists volumes (GET /api/service/:id/volumes).
func (c *Client) ListAIVolumes(ctx context.Context, serviceID string) ([]AIVolume, error) {
	var raw any
	path := "/api/service/" + PathEscape(serviceID) + "/volumes"
	if err := c.Get(ctx, path, nil, &raw); err != nil {
		return nil, err
	}
	list := ExtractList(raw, "volumes")
	out := make([]AIVolume, 0, len(list))
	for _, item := range list {
		if m, ok := AsMap(item); ok {
			out = append(out, parseAIVolume(m))
		}
	}
	return out, nil
}

// GetAIVolume fetches one volume (GET /api/service/:id/volumes/:volume_id).
func (c *Client) GetAIVolume(ctx context.Context, serviceID, volumeID string) (*AIVolume, error) {
	var raw any
	path := "/api/service/" + PathEscape(serviceID) + "/volumes/" + PathEscape(volumeID)
	if err := c.Get(ctx, path, nil, &raw); err != nil {
		return nil, err
	}
	m, ok := AsMap(raw)
	if !ok {
		if list := ExtractList(raw, "volumes"); len(list) > 0 {
			m, _ = AsMap(list[0])
		}
	}
	if m == nil {
		return nil, fmt.Errorf("unexpected response shape for volume %s", volumeID)
	}
	vol := parseAIVolume(m)
	if vol.ID == "" {
		vol.ID = volumeID
	}
	return &vol, nil
}

// CreateAIVolume provisions a volume (POST /api/service/:id/volumes).
func (c *Client) CreateAIVolume(ctx context.Context, serviceID string, in AIVolume) (string, error) {
	var raw any
	path := "/api/service/" + PathEscape(serviceID) + "/volumes"
	q := Query("cloud", in.Cloud, "region", in.Region, "name", in.Name)
	if in.SizeInGB > 0 {
		q.Set("size_in_gb", fmt.Sprintf("%d", in.SizeInGB))
	}
	if err := c.Post(ctx, path, q, &raw); err != nil {
		return "", err
	}
	if id := FirstString(raw, "id", "volume_id", "volumeid"); id != "" {
		return id, nil
	}
	// Fallback: locate by name.
	volumes, err := c.ListAIVolumes(ctx, serviceID)
	if err != nil {
		return "", fmt.Errorf("volume was created but its id could not be resolved: %w", err)
	}
	for _, vol := range volumes {
		if in.Name != "" && vol.Name == in.Name {
			return vol.ID, nil
		}
	}
	return "", fmt.Errorf("volume %q was created but could not be found", in.Name)
}

// DeleteAIVolume removes a volume (POST /api/service/:id/volumes/:volume_id/delete).
func (c *Client) DeleteAIVolume(ctx context.Context, serviceID, volumeID string) error {
	path := "/api/service/" + PathEscape(serviceID) + "/volumes/" + PathEscape(volumeID) + "/delete"
	return c.Post(ctx, path, nil, nil)
}

func parseAIVolume(m map[string]any) AIVolume {
	return AIVolume{
		ID:       FirstString(m, "id", "volume_id", "volumeid"),
		Name:     FirstString(m, "name", "label", "volume_name"),
		Cloud:    FirstString(m, "cloud", "provider"),
		Region:   FirstString(m, "region", "location"),
		SizeInGB: FirstInt64(m, "size_in_gb", "size", "size_gb", "gb"),
		Status:   FirstString(m, "status", "state"),
	}
}

// ---------------------------------------------------------------------------
// Templates
// ---------------------------------------------------------------------------

// AITemplate is a saved instance template.
type AITemplate struct {
	ID          string
	Name        string
	Description string
	IsPublic    bool
}

// ListAITemplates lists templates (GET /api/service/:id/templates).
func (c *Client) ListAITemplates(ctx context.Context, serviceID string) ([]AITemplate, error) {
	var raw any
	path := "/api/service/" + PathEscape(serviceID) + "/templates"
	if err := c.Get(ctx, path, nil, &raw); err != nil {
		return nil, err
	}
	list := ExtractList(raw, "templates")
	out := make([]AITemplate, 0, len(list))
	for _, item := range list {
		if m, ok := AsMap(item); ok {
			out = append(out, parseAITemplate(m))
		}
	}
	return out, nil
}

// GetAITemplate fetches one template (GET /api/service/:id/templates/:template_id).
func (c *Client) GetAITemplate(ctx context.Context, serviceID, templateID string) (*AITemplate, error) {
	var raw any
	path := "/api/service/" + PathEscape(serviceID) + "/templates/" + PathEscape(templateID)
	if err := c.Get(ctx, path, nil, &raw); err != nil {
		return nil, err
	}
	m, ok := AsMap(raw)
	if !ok {
		if list := ExtractList(raw, "templates"); len(list) > 0 {
			m, _ = AsMap(list[0])
		}
	}
	if m == nil {
		return nil, fmt.Errorf("unexpected response shape for template %s", templateID)
	}
	tpl := parseAITemplate(m)
	if tpl.ID == "" {
		tpl.ID = templateID
	}
	return &tpl, nil
}

// CreateAITemplate saves a template (POST /api/service/:id/templates).
func (c *Client) CreateAITemplate(ctx context.Context, serviceID string, in AITemplate) (string, error) {
	var raw any
	path := "/api/service/" + PathEscape(serviceID) + "/templates"
	q := Query("name", in.Name, "description", in.Description)
	q.Set("is_public", fmt.Sprintf("%t", in.IsPublic))
	if err := c.Post(ctx, path, q, &raw); err != nil {
		return "", err
	}
	if id := FirstString(raw, "id", "template_id", "templateid"); id != "" {
		return id, nil
	}
	// Fallback: locate by name.
	templates, err := c.ListAITemplates(ctx, serviceID)
	if err != nil {
		return "", fmt.Errorf("template was created but its id could not be resolved: %w", err)
	}
	for _, tpl := range templates {
		if tpl.Name == in.Name {
			return tpl.ID, nil
		}
	}
	return "", fmt.Errorf("template %q was created but could not be found", in.Name)
}

// UpdateAITemplate updates a template
// (POST /api/service/:id/templates/:template_id/update).
func (c *Client) UpdateAITemplate(ctx context.Context, serviceID, templateID string, in AITemplate) error {
	path := "/api/service/" + PathEscape(serviceID) + "/templates/" + PathEscape(templateID) + "/update"
	q := Query("name", in.Name, "description", in.Description)
	q.Set("is_public", fmt.Sprintf("%t", in.IsPublic))
	return c.Post(ctx, path, q, nil)
}

// DeleteAITemplate removes a template
// (POST /api/service/:id/templates/:template_id/delete).
func (c *Client) DeleteAITemplate(ctx context.Context, serviceID, templateID string) error {
	path := "/api/service/" + PathEscape(serviceID) + "/templates/" + PathEscape(templateID) + "/delete"
	return c.Post(ctx, path, nil, nil)
}

func parseAITemplate(m map[string]any) AITemplate {
	return AITemplate{
		ID:          FirstString(m, "id", "template_id", "templateid"),
		Name:        FirstString(m, "name", "label", "template_name"),
		Description: FirstString(m, "description", "desc"),
		IsPublic:    FirstBool(m, "is_public", "public", "ispublic"),
	}
}

// ---------------------------------------------------------------------------
// Clusters
// ---------------------------------------------------------------------------

// AICluster is a GPU cluster.
type AICluster struct {
	ID           string
	Name         string
	Cloud        string
	Region       string
	ClusterType  string
	NumInstances int64
	SSHKeyID     string
	OS           string
	Status       string
}

// ListAIClusters lists clusters (GET /api/service/:id/clusters).
func (c *Client) ListAIClusters(ctx context.Context, serviceID string) ([]AICluster, error) {
	var raw any
	path := "/api/service/" + PathEscape(serviceID) + "/clusters"
	if err := c.Get(ctx, path, nil, &raw); err != nil {
		return nil, err
	}
	list := ExtractList(raw, "clusters")
	out := make([]AICluster, 0, len(list))
	for _, item := range list {
		if m, ok := AsMap(item); ok {
			out = append(out, parseAICluster(m))
		}
	}
	return out, nil
}

// GetAICluster fetches one cluster (GET /api/service/:id/clusters/:cluster_id).
func (c *Client) GetAICluster(ctx context.Context, serviceID, clusterID string) (*AICluster, error) {
	var raw any
	path := "/api/service/" + PathEscape(serviceID) + "/clusters/" + PathEscape(clusterID)
	if err := c.Get(ctx, path, nil, &raw); err != nil {
		return nil, err
	}
	m, ok := AsMap(raw)
	if !ok {
		if list := ExtractList(raw, "clusters"); len(list) > 0 {
			m, _ = AsMap(list[0])
		}
	}
	if m == nil {
		return nil, fmt.Errorf("unexpected response shape for cluster %s", clusterID)
	}
	cl := parseAICluster(m)
	if cl.ID == "" {
		cl.ID = clusterID
	}
	return &cl, nil
}

// CreateAICluster provisions a GPU cluster (POST /api/service/:id/clusters).
func (c *Client) CreateAICluster(ctx context.Context, serviceID string, in AICluster) (string, error) {
	var raw any
	path := "/api/service/" + PathEscape(serviceID) + "/clusters"
	q := Query(
		"cloud", in.Cloud,
		"region", in.Region,
		"cluster_type", in.ClusterType,
		"name", in.Name,
		"ssh_key_id", in.SSHKeyID,
		"os", in.OS,
	)
	if in.NumInstances > 0 {
		q.Set("num_instances", fmt.Sprintf("%d", in.NumInstances))
	}
	if err := c.Post(ctx, path, q, &raw); err != nil {
		return "", err
	}
	if id := FirstString(raw, "id", "cluster_id", "clusterid"); id != "" {
		return id, nil
	}
	// Fallback: locate by name.
	clusters, err := c.ListAIClusters(ctx, serviceID)
	if err != nil {
		return "", fmt.Errorf("cluster was created but its id could not be resolved: %w", err)
	}
	for _, cl := range clusters {
		if cl.Name == in.Name {
			return cl.ID, nil
		}
	}
	return "", fmt.Errorf("cluster %q was created but could not be found", in.Name)
}

// DeleteAICluster removes a cluster
// (POST /api/service/:id/clusters/:cluster_id/delete).
func (c *Client) DeleteAICluster(ctx context.Context, serviceID, clusterID string) error {
	path := "/api/service/" + PathEscape(serviceID) + "/clusters/" + PathEscape(clusterID) + "/delete"
	return c.Post(ctx, path, nil, nil)
}

func parseAICluster(m map[string]any) AICluster {
	return AICluster{
		ID:           FirstString(m, "id", "cluster_id", "clusterid"),
		Name:         FirstString(m, "name", "label", "cluster_name"),
		Cloud:        FirstString(m, "cloud", "provider"),
		Region:       FirstString(m, "region", "location"),
		ClusterType:  FirstString(m, "cluster_type", "type", "cluster_type_id"),
		NumInstances: FirstInt64(m, "num_instances", "instances", "nodes", "num_nodes"),
		SSHKeyID:     asString(First(m, "ssh_key_id", "sshkey_id")),
		OS:           FirstString(m, "os", "image"),
		Status:       FirstString(m, "status", "state"),
	}
}

// ---------------------------------------------------------------------------
// Catalog lookups (instance types, volume types, cluster types)
// ---------------------------------------------------------------------------

// AIInstanceType is an entry of GET /api/service/:id/instances/types.
type AIInstanceType struct {
	ID    string
	Name  string
	GPU   string
	VCPUs int64
	RAM   string
}

// ListAIInstanceTypes lists available GPU instance types.
func (c *Client) ListAIInstanceTypes(ctx context.Context, serviceID string) ([]AIInstanceType, error) {
	var raw any
	path := "/api/service/" + PathEscape(serviceID) + "/instances/types"
	if err := c.Get(ctx, path, nil, &raw); err != nil {
		return nil, err
	}
	list := ExtractList(raw, "types", "instance_types", "instanceTypes")
	out := make([]AIInstanceType, 0, len(list))
	for _, item := range list {
		if m, ok := AsMap(item); ok {
			out = append(out, AIInstanceType{
				ID:    FirstString(m, "id", "type", "name", "instance_type"),
				Name:  FirstString(m, "name", "display_name", "label"),
				GPU:   FirstString(m, "gpu", "gpu_type", "gpu_name"),
				VCPUs: FirstInt64(m, "vcpus", "cpu", "cpus", "vcpu"),
				RAM:   FirstString(m, "ram", "memory", "memory_gb"),
			})
		} else {
			// The catalog may be a plain array of type names.
			if s := asString(item); s != "" {
				out = append(out, AIInstanceType{ID: s, Name: s})
			}
		}
	}
	return out, nil
}

// AIVolumeType is an entry of GET /api/service/:id/volumes/types.
type AIVolumeType struct {
	ID   string
	Name string
}

// ListAIVolumeTypes lists available volume types.
func (c *Client) ListAIVolumeTypes(ctx context.Context, serviceID string) ([]AIVolumeType, error) {
	var raw any
	path := "/api/service/" + PathEscape(serviceID) + "/volumes/types"
	if err := c.Get(ctx, path, nil, &raw); err != nil {
		return nil, err
	}
	list := ExtractList(raw, "types", "volume_types", "volumeTypes")
	out := make([]AIVolumeType, 0, len(list))
	for _, item := range list {
		if m, ok := AsMap(item); ok {
			out = append(out, AIVolumeType{
				ID:   FirstString(m, "id", "type", "name"),
				Name: FirstString(m, "name", "display_name", "label"),
			})
		} else {
			if s := asString(item); s != "" {
				out = append(out, AIVolumeType{ID: s, Name: s})
			}
		}
	}
	return out, nil
}

// AIClusterType is an entry of GET /api/service/:id/clusters/types.
type AIClusterType struct {
	ID   string
	Name string
}

// ListAIClusterTypes lists available GPU cluster types.
func (c *Client) ListAIClusterTypes(ctx context.Context, serviceID string) ([]AIClusterType, error) {
	var raw any
	path := "/api/service/" + PathEscape(serviceID) + "/clusters/types"
	if err := c.Get(ctx, path, nil, &raw); err != nil {
		return nil, err
	}
	list := ExtractList(raw, "types", "cluster_types", "clusterTypes")
	out := make([]AIClusterType, 0, len(list))
	for _, item := range list {
		if m, ok := AsMap(item); ok {
			out = append(out, AIClusterType{
				ID:   FirstString(m, "id", "type", "name"),
				Name: FirstString(m, "name", "display_name", "label"),
			})
		} else {
			if s := asString(item); s != "" {
				out = append(out, AIClusterType{ID: s, Name: s})
			}
		}
	}
	return out, nil
}
