[OpenAPI specs](./bundled.yaml) is copied from [Kibana repo](https://raw.githubusercontent.com/elastic/kibana/8.7/x-pack/plugins/actions/docs/openapi/bundled.yaml) with some modifications:

- added mapping section for discriminator field in `POST` `/s/{spaceId}/api/actions/connector`;
- added explicit object definitions for `400`, `401` and `404` errors (`oapi-codegen` doesn't generate proper code for embedded anonymous objects in some cases) - `bad_request_error`, `authorization_error` and `object_not_found_error`;
- added missing `oneOf` types in `requestBody` for `PUT` `/s/{spaceId}/api/actions/connector/{connectorId}` - the original `bundled.yaml` misses some connector types in the `PUT` `requestBody` defintion:
  - `update_connector_request_email`;
  - `update_connector_request_pagerduty`;
  - `update_connector_request_servicenow_sir`;
  - `update_connector_request_slack`;
  - `update_connector_request_teams`;
  - `update_connector_request_tines`;
  - `update_connector_request_webhook`;
  - `update_connector_request_xmatters`.
- response definitions of `/s/{spaceId}/api/actions/connector/{connectorId}/_execute` and `/s/{spaceId}/api/actions/action/{actionId}/_execute` are modified from embedded object definitions to named ones `run_connector_general_response` and `legacy_run_connector_general_response`;
- specified properties for following types. The original `bundled.yaml` defines them as dynamic objects (`additionalProperties: true`):
  - `config_propeties_email`;
  - `config_properties_pagerduty`;
  - `config_properties_tines`;
  - `config_properties_webhook`;
  - `config_properties_xmatters`;
- `is_deprecated` is marked as optional field (it's required field in the vanilla `bundled.yaml`) in the following objects (Kibana responses may omit it):
  - `connector_response_properties_cases_webhook`; 
  - `connector_response_properties_email`;
  - `connector_response_properties_index`;
  - `connector_response_properties_jira`;
  - `connector_response_properties_opsgenie`;
  - `connector_response_properties_pagerduty`;
  - `connector_response_properties_resilient`;
  - `connector_response_properties_serverlog`;
  - `connector_response_properties_servicenow`;
  - `connector_response_properties_servicenow_itom`;
  - `connector_response_properties_servicenow_sir`;
  - `connector_response_properties_slack`;
  - `connector_response_properties_swimlane`;
  - `connector_response_properties_teams`;
  - `connector_response_properties_tines`;
  - `connector_response_properties_webhook`;
  - `connector_response_properties_xmatters`.
- added mapping section for discriminator field in `connector_response_properties`.