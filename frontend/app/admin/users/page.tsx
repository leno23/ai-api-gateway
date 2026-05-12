"use client";

import { App, Button, Card, Form, InputNumber, Popconfirm, Typography } from "antd";
import { useState } from "react";

import { ApiError } from "@/lib/api/client";
import { useGatewayRequest } from "@/lib/hooks/use-gateway-request";

export default function AdminUsersPage() {
  const request = useGatewayRequest();
  const { message } = App.useApp();
  const [form] = Form.useForm<{ user_id: number; status: number }>();
  const [loading, setLoading] = useState(false);

  return (
    <div className="flex max-w-lg flex-col gap-4">
      <div>
        <Typography.Title level={4} className="!mb-1">
          用户治理
        </Typography.Title>
        <Typography.Paragraph type="secondary" className="!mb-0 text-sm">
          MVP：按用户 ID 更新状态（1 正常 / 0 封禁）。用户列表与搜索依赖后端扩展接口。
        </Typography.Paragraph>
      </div>
      <Card size="small">
        <Form
          form={form}
          layout="vertical"
          initialValues={{ status: 0 }}
          onFinish={async (v) => {
            setLoading(true);
            try {
              await request(`/admin/users/${v.user_id}/status`, {
                method: "PATCH",
                body: JSON.stringify({ status: v.status }),
              });
              message.success("已更新用户状态");
              form.resetFields(["user_id"]);
            } catch (e) {
              if (e instanceof ApiError) {
                message.error(e.message || "更新失败");
              }
            } finally {
              setLoading(false);
            }
          }}
        >
          <Form.Item
            name="user_id"
            label="用户 ID"
            rules={[{ required: true, type: "number", min: 1, message: "请输入有效用户 ID" }]}
          >
            <InputNumber className="w-full" placeholder="整数 ID" />
          </Form.Item>
          <Form.Item
            name="status"
            label="状态"
            rules={[{ required: true, type: "number" }]}
            extra="1 = 正常，0 = 禁用"
          >
            <InputNumber className="w-full" min={0} max={1} />
          </Form.Item>
          <Form.Item>
            <Popconfirm
              title="确认修改该用户状态？"
              onConfirm={() => form.submit()}
              okText="确认"
              cancelText="取消"
            >
              <Button type="primary" loading={loading}>
                提交
              </Button>
            </Popconfirm>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
}
