drop table if exists operation_logs;

CREATE TABLE operation_logs(
    id             SERIAL PRIMARY KEY,
    operation_type TEXT NOT NULL,
    name           TEXT NOT NULL,
    record_id      INTEGER NOT NULL,
    ip_address     TEXT NOT NULL,
    user_agent     TEXT NOT NULL,
    created_at     TIMESTAMP DEFAULT NOW() NOT NULL
);;

COMMENT ON COLUMN operation_logs.id IS '日志ID';
COMMENT ON COLUMN operation_logs.operation_type IS '操作类型';
COMMENT ON COLUMN operation_logs.name IS '表名';
COMMENT ON COLUMN operation_logs.record_id IS '记录ID';
COMMENT ON COLUMN operation_logs.ip_address IS 'IP地址';
COMMENT ON COLUMN operation_logs.user_agent IS '用户代理';
COMMENT ON COLUMN operation_logs.created_at IS '创建时间';
