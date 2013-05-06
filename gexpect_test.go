package gexpect

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"testing"
)

func TestPlan(t *testing.T) {
	sshPlan := NewPlan()
	loginStep := sshPlan.Flow.Expect(5, "[Uu]sername:")
	loginStep.WhenTimeout.Terminate()
	loginStep.WhenMatched.VarSendLine("username")
	passStep := loginStep.WhenMatched.Expect(5, "[Pp]assword:")
	passStep.WhenTimeout.Terminate()
	passStep.WhenMatched.VarSendLine("password")
	passStep.WhenMatched.Expect(5, "[$#>]")
	/*fn := func(step *Step) error {
		spew.Dump(step)
		return nil
	}*/
	//Walk(sshPlan.Flow, fn)
	sshPlan.BindVar("username", "knightmare")
	sshPlan.BindVar("password", "pass")
	if err := sshPlan.BindApply(); err != nil {
		fmt.Println(err)
	}
	//Walk(sshPlan.Flow, fn)
	spew.Dump(sshPlan)
	//shellPrompt.interactive()
}
