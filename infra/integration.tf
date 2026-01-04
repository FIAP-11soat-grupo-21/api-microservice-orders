resource "aws_lb_listener" "listener" {
  load_balancer_arn = data.terraform_remote_state.infra.outputs.alb_arn
  port              = var.image_port
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = data.terraform_remote_state.infra.outputs.alb_target_group_arn
  }

  tags = merge(
    data.terraform_remote_state.infra.outputs.project_common_tags
    , { Name = "${var.application_name}-listener" }
  )
}

resource "aws_alb_listener_rule" "rule" {
  listener_arn = aws_lb_listener.listener.arn

  condition {
    path_pattern {
      values = var.app_path_pattern
    }
  }

  action {
    type             = "forward"
    target_group_arn = data.terraform_remote_state.infra.outputs.alb_target_group_arn
  }
}

resource "aws_apigatewayv2_integration" "alb_proxy" {
  api_id           = data.terraform_remote_state.infra.outputs.api_gateway_id
  integration_type = var.apigw_integration_type

  integration_uri        = aws_lb_listener.listener.arn
  integration_method     = var.apigw_integration_method
  payload_format_version = var.apigw_payload_format_version

  connection_type = var.apigw_connection_type
  connection_id   = data.terraform_remote_state.infra.outputs.api_gateway_vpc_link_id
}

