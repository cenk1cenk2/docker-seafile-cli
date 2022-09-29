package pipe

type Ctx struct {
	Health struct {
		SeafDaemonPID int
		CcnetPID      int
	}

	Libraries []string
}
