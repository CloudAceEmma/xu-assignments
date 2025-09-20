# Todo 应用后端 API

基于golang gin framework 的 RESTful 待办事项管理 API，支持完整的 CRUD、状态切换、批量删除、分页、过滤、健康检查、Swagger 文档等。

## 技术栈
- golang
- go-gin
- mariadb 12.0.2+
- Swagger

## 数据库设计

- 地址：localhost:3306
- 用户名：monty
- 密码：test001
- 数据库名：todoapp

## API接口规范

### 基础信息
- **基础URL**: `http://localhost:8000`
- **API前缀**: `/api/v1`
- **数据格式**: JSON
- **HTTP状态码**: 遵循RESTful规范

### 接口列表

#### 1. 获取所有待办事项
```http
GET /api/v1/todos
```

**查询参数**:
- `completed` (可选): boolean - 过滤完成状态
- `limit` (可选): integer - 限制返回数量，默认100
- `offset` (可选): integer - 偏移量，默认0

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "title": "学习React",
      "description": "完成React基础教程",
      "completed": false,
      "created_at": "2024-01-01T10:00:00Z",
      "updated_at": "2024-01-01T10:00:00Z",
      "priority": 1,
      "due_date": null
    }
  ],
  "total": 1
}
```

#### 2. 创建待办事项
```http
POST /api/v1/todos
```

**请求体**:
```json
{
  "title": "新的待办事项",
  "description": "描述信息（可选）",
  "priority": 0,
  "due_date": null
}
```

**响应示例**:
```json
{
  "code": 201,
  "message": "Todo created successfully",
  "data": {
    "id": 2,
    "title": "新的待办事项",
    "description": "描述信息（可选）",
    "completed": false,
    "created_at": "2024-01-01T11:00:00Z",
    "updated_at": "2024-01-01T11:00:00Z",
    "priority": 0,
    "due_date": null
  }
}
```

#### 3. 更新待办事项
```http
PUT /api/v1/todos/{todo_id}
```

**路径参数**:
- `todo_id`: integer - 待办事项ID

**请求体**:
```json
{
  "title": "更新的标题",
  "description": "更新的描述",
  "completed": true,
  "priority": 2,
  "due_date": "2024-12-31T23:59:59Z"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "Todo updated successfully",
  "data": {
    "id": 1,
    "title": "更新的标题",
    "description": "更新的描述",
    "completed": true,
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z",
    "priority": 2,
    "due_date": "2024-12-31T23:59:59Z"
  }
}
```

#### 4. 删除待办事项
```http
DELETE /api/v1/todos/{todo_id}
```

**路径参数**:
- `todo_id`: integer - 待办事项ID

**响应示例**:
```json
{
  "code": 200,
  "message": "Todo deleted successfully"
}
```

#### 5. 标记完成/取消完成
```http
PATCH /api/v1/todos/{todo_id}/toggle
```

**路径参数**:
- `todo_id`: integer - 待办事项ID

**响应示例**:
```json
{
  "code": 200,
  "message": "Todo status toggled successfully",
  "data": {
    "id": 1,
    "completed": true,
    "updated_at": "2024-01-01T13:00:00Z"
  }
}
```

#### 6. 批量删除已完成项
```http
DELETE /api/v1/todos/completed
```

**响应示例**:
```json
{
  "code": 200,
  "message": "Completed todos deleted successfully",
  "deleted_count": 5
}
```

#### 7. 清空所有待办事项
```http
DELETE /api/v1/todos/all
```

**响应示例**:
```json
{
  "code": 200,
  "message": "All todos deleted successfully",
  "deleted_count": 10
}
```

### 错误响应格式

```json
{
  "code": 404,
  "message": "Todo not found",
  "detail": "Todo with id 999 does not exist"
}
```






## 测试策略

1. **单元测试**: 
2. **集成测试**: API接口测试
3. **E2E测试**: 
4. **性能测试**: 大量数据场景测试

