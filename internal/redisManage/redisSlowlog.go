package redisManage

func (ro *RedisOps) SlowLog() error {
	//client, err := NewRedisGo()
	//if err != nil {
	//	return errors.WithStack(err)
	//}
	//defer client.Close()
	//slowlogs, err := client.SlowLog()
	//if err != nil {
	//	return errors.WithStack(err)
	//}
	//headers := []string{"ID", "开始时间", "耗时(秒)", "命令", "客户端地址"}
	//data := make([][]string, 0, len(slowlogs))
	//for _, v := range slowlogs {
	//	idStr := strconv.FormatInt(v.ID, 10)
	//	startTime := v.Time.String()
	//	executionTime := fmt.Sprintf("%f", v.ExecutionTime.Seconds())
	//	args := fmt.Sprintf("%s", v.Args)
	//	data = append(data, []string{idStr, startTime, executionTime, args, v.ClientAddr})
	//}
	//tableRender.Render(headers, data)
	return nil
}
