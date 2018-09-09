package global

const (
	// AppStatus 记录app的status
	AppStatus = "app-status"
)

const (
	// AppRunning app's running
	AppRunning = iota + 1
	// AppPause app's pause
	AppPause
	// AppStop app's stop
	AppStop
)

// StartApplication 设置应用状态为running
func StartApplication() {
	Store(AppStatus, AppRunning)
}

// PauseApplication 设置应用状态为pause
func PauseApplication() {
	Store(AppStatus, AppPause)
}

// IsApplicationRunning 判断是否正在运行
func IsApplicationRunning() bool {
	v, ok := Load(AppStatus)
	if !ok {
		return false
	}
	return v.(int) == AppRunning
}
