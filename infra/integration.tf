resource "aws_alb_target_group" "target_group" {
  name        = "${var.application_name}-tg"
  port        = var.image_port
  protocol    = "HTTP"
  vpc_id      = data.terraform_remote_state.infra.outputs.vpc_id
  target_type = "ip"

  health_check {
    path                = var.health_check_path
    protocol            = "HTTP"
    matcher             = "200-399"
    interval            = 30
    timeout             = 5
    healthy_threshold   = 2
    unhealthy_threshold = 2
  }

  tags = merge(
    data.terraform_remote_state.infra.outputs.project_common_tags
    , { Name = "${var.application_name}-target-group" }
  )
}

resource "aws_lb_listener" "listener" {
  depends_on = [aws_alb_target_group.target_group]

  load_balancer_arn = data.terraform_remote_state.infra.outputs.alb_arn
  port              = var.image_port
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_alb_target_group.target_group.arn
  }

  tags = merge(
    data.terraform_remote_state.infra.outputs.project_common_tags
    , { Name = "${var.application_name}-listener" }
  )
}

resource "aws_alb_listener_rule" "rule" {
  depends_on = [aws_lb_listener.listener, aws_alb_target_group.target_group]

  listener_arn = aws_lb_listener.listener.arn
  condition {
    path_pattern {
      values = var.app_path_pattern
    }
  }
  action {
    type             = "forward"
    target_group_arn = aws_alb_target_group.target_group.arn
  }
}

resource "aws_apigatewayv2_integration" "alb_proxy" {
  api_id                 = data.terraform_remote_state.infra.outputs.api_gateway_id
  integration_type       = var.apigw_integration_type
  integration_uri        = aws_lb_listener.listener.arn
  integration_method     = var.apigw_integration_method
  payload_format_version = var.apigw_payload_format_version

  connection_type = var.apigw_connection_type
  connection_id   = data.terraform_remote_state.infra.outputs.api_gateway_vpc_link_id
}

