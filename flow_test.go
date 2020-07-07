package workflow_test

import (
	"encoding/json"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmniu/workflow"
	"github.com/jmniu/workflow/service/db"
)

func init() {
	var err error
	workflow.Init(
		db.SetDSN("root:GmTech@2019@tcp(192.168.238.178:3306)/db_flow?charset=utf8"),
		db.SetTrace(false),
	)

	//err = workflow.LoadFile("test_data/leave.bpmn")
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = workflow.LoadFile("test_data/apply_sqltest.bpmn")
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = workflow.LoadFile("test_data/parallel_test.bpmn")
	//if err != nil {
	//	panic(err)
	//}

	//err = workflow.LoadFile("test_data/form_test.bpmn")
	//if err != nil {
	//	panic(err)
	//}

	err = workflow.LoadFile("test_data/proc_form.bpmn")
	if err != nil {
		panic(err)
	}
}

func TestRepair(t *testing.T) {
	var flowCode = "id_process_repair"
	var input = map[string]interface{}{
		"repair": "niujiaming",
	}
	result, err := workflow.StartFlow(flowCode, "id_start", "niujiaming", input)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Printf("%v", result)

	input["verify"] = "niujiaming"
	result, err = workflow.HandleFlow(result.NextNodes[0].NodeInstance.RecordID, "niujiaming", input)
	if err != nil {
		t.Fatal(err.Error())
	}

	input["ok"] = true
	result, err = workflow.HandleFlow(result.NextNodes[0].NodeInstance.RecordID, "niujiaming", input)
	if err != nil {
		t.Fatal(err.Error())
	}

	fmt.Printf("%v", result)
}

func TestLeaveBzrApprovalPass(t *testing.T) {
	var (
		flowCode = "process_leave_test"
		bzr      = "T002"
	)

	input := map[string]interface{}{
		"day": 1,
		"bzr": bzr,
	}

	// 开始流程
	result, err := workflow.StartFlow(flowCode, "node_start", "T001", input)
	if err != nil {
		t.Fatal(err.Error())
	}

	if result.NextNodes[0].CandidateIDs[0] != bzr {
		t.Fatalf("无效的下一级流转：%s", result.String())
	}

	// 查询待办
	todos, err := workflow.QueryTodoFlows(flowCode, bzr)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(todos) != 1 {
		bts, _ := json.Marshal(todos)
		t.Fatalf("无效的待办数据:%s", string(bts))
	}

	// 处理流程（通过）
	input["action"] = "pass"
	result, err = workflow.HandleFlow(todos[0].RecordID, bzr, input)
	if err != nil {
		t.Fatal(err.Error())
	}

	// 流程结束
	if !result.IsEnd {
		t.Fatalf("无效的处理结果：%s", result.String())
	}
}

func TestLeaveBzrApprovalBack(t *testing.T) {
	var (
		flowCode = "process_leave_test"
		launcher = "T001"
		bzr      = "T002"
	)

	input := map[string]interface{}{
		"day": 1,
		"bzr": bzr,
	}

	// 开始流程
	result, err := workflow.StartFlow(flowCode, "node_start", launcher, input)
	if err != nil {
		t.Fatal(err.Error())
	}

	if result.NextNodes[0].CandidateIDs[0] != bzr {
		t.Fatalf("无效的下一级流转：%s", result.String())
	}

	// 查询待办
	todos, err := workflow.QueryTodoFlows(flowCode, bzr)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// 处理流程（退回）
	input["action"] = "back"
	result, err = workflow.HandleFlow(todos[0].RecordID, bzr, input)
	if err != nil {
		t.Fatal(err.Error())
	}

	if result.IsEnd ||
		result.NextNodes[0].CandidateIDs[0] != launcher {
		t.Fatalf("无效的处理结果：%s", result.String())
	}

	// 查询退回流程
	todos, err = workflow.QueryTodoFlows(flowCode, launcher)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// 处理退回流程
	delete(input, "action")
	result, err = workflow.HandleFlow(todos[0].RecordID, launcher, input)
	if err != nil {
		t.Fatal(err.Error())
	}

	if result.NextNodes[0].CandidateIDs[0] != bzr {
		t.Fatalf("无效的下一级流转：%s", result.String())
	}

	// 查询待办流程
	todos, err = workflow.QueryTodoFlows(flowCode, bzr)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// 处理流程（通过）
	input["action"] = "pass"
	result, err = workflow.HandleFlow(todos[0].RecordID, bzr, input)
	if err != nil {
		t.Fatal(err.Error())
	}

	// 流程结束
	if !result.IsEnd {
		t.Fatalf("无效的处理结果：%s", result.String())
	}
}

func TestLeaveFdyApprovalPass(t *testing.T) {
	var (
		flowCode = "process_leave_test"
		launcher = "T001"
		bzr      = "T002"
		fdy      = "T003"
	)

	input := map[string]interface{}{
		"day": 3,
		"bzr": bzr,
		"fdy": fdy,
	}

	// 开始流程
	result, err := workflow.StartFlow(flowCode, "node_start", launcher, input)
	if err != nil {
		t.Fatal(err.Error())
	}

	if result.NextNodes[0].CandidateIDs[0] != bzr {
		t.Fatalf("无效的下一级流转：%s", result.String())
	}

	// 查询待办
	todos, err := workflow.QueryTodoFlows(flowCode, bzr)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// 处理流程（通过）
	input["action"] = "pass"
	result, err = workflow.HandleFlow(todos[0].RecordID, bzr, input)
	if err != nil {
		t.Fatal(err.Error())
	}

	// 查询待办
	todos, err = workflow.QueryTodoFlows(flowCode, fdy)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// 处理流程（通过）
	input["action"] = "pass"
	result, err = workflow.HandleFlow(todos[0].RecordID, fdy, input)
	if err != nil {
		t.Fatal(err.Error())
	}

	// 流程结束
	if !result.IsEnd {
		t.Fatalf("无效的处理结果：%s", result.String())
	}
}

func TestApplySQLPass(t *testing.T) {
	var (
		flowCode = "process_apply_sqltest"
	)

	input := map[string]interface{}{
		"form": "apply",
	}

	// 开始流程
	result, err := workflow.StartFlow(flowCode, "node_start", "A001", input)
	if err != nil {
		t.Fatal(err.Error())
	}

	cIDs := result.NextNodes[0].CandidateIDs
	if len(cIDs) != 2 {
		t.Fatalf("无效的下一级流转：%s", result.String())
	}

	var (
		nodeInstanceID string
		userID         string
	)
	for _, cid := range cIDs {
		// 查询待办
		todos, err := workflow.QueryTodoFlows(flowCode, cid)
		if err != nil {
			t.Fatalf(err.Error())
		}

		if len(todos) != 1 {
			bts, _ := json.Marshal(todos)
			t.Fatalf("无效的待办数据:%s", string(bts))
		}

		nodeInstanceID = todos[0].RecordID
		userID = cid
	}

	// 处理流程（通过）
	input["action"] = "pass"
	result, err = workflow.HandleFlow(nodeInstanceID, userID, input)
	if err != nil {
		t.Fatal(err.Error())
	}

	// 流程结束
	if !result.IsEnd {
		t.Fatalf("无效的处理结果：%s", result.String())
	}
}

func TestParallel(t *testing.T) {
	var (
		flowCode = "process_parallel_test"
	)

	input := map[string]interface{}{
		"form": "countersign",
	}

	// 开始流程
	result, err := workflow.StartFlow(flowCode, "node_start", "H001", input)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(result.NextNodes) != 3 {
		t.Fatalf("无效的下一级流转：%s", result.String())
	}

	for i, node := range result.NextNodes {
		if len(node.CandidateIDs) != 1 {
			t.Fatalf("无效的节点处理人：%v", node.CandidateIDs)
		}

		todos, err := workflow.QueryTodoFlows(flowCode, node.CandidateIDs[0])
		if err != nil {
			t.Fatalf(err.Error())
		} else if len(todos) != 1 {
			bts, _ := json.Marshal(todos)
			t.Fatalf("无效的待办数据:%s", string(bts))
		}

		input["sign"] = node.CandidateIDs[0]
		result, err := workflow.HandleFlow(todos[0].RecordID, node.CandidateIDs[0], input)
		if err != nil {
			t.Fatalf(err.Error())
		}

		if i == 2 {
			if !result.IsEnd {
				t.Fatalf("无效的处理结果：%s", result.String())
			}
			break
		}

		if result.IsEnd {
			t.Fatalf("无效的处理结果：%s", result.String())
		}
	}

}

func TestLeaveRepeatedBack(t *testing.T) {
	var (
		flowCode = "process_leave_test"
		launcher = "B001"
		bzr      = "B002"
		fdy      = "B003"
		yld      = "B004"
	)

	input := map[string]interface{}{
		"day": 5,
		"bzr": bzr,
		"fdy": fdy,
		"yld": yld,
	}

	// 开始流程
	result, err := workflow.StartFlow(flowCode, "node_start", launcher, input)
	if err != nil {
		t.Fatal(err.Error())
	}

	if result.NextNodes[0].CandidateIDs[0] != bzr {
		t.Fatalf("无效的下一级流转：%s", result.String())
	}

	// 查询待办
	todos, err := workflow.QueryTodoFlows(flowCode, bzr)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// 处理流程（通过）
	input["action"] = "pass"
	result, err = workflow.HandleFlow(todos[0].RecordID, bzr, input)
	if err != nil {
		t.Fatal(err.Error())
	}

	// 查询待办
	todos, err = workflow.QueryTodoFlows(flowCode, fdy)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// 处理流程（通过）
	input["action"] = "pass"
	result, err = workflow.HandleFlow(todos[0].RecordID, fdy, input)
	if err != nil {
		t.Fatal(err.Error())
	}

	// 查询待办
	todos, err = workflow.QueryTodoFlows(flowCode, yld)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// 处理流程（退回）
	input["action"] = "back"
	result, err = workflow.HandleFlow(todos[0].RecordID, fdy, input)
	if err != nil {
		t.Fatal(err.Error())
	}

	// 查询待办
	todos, err = workflow.QueryTodoFlows(flowCode, launcher)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// 处理流程
	result, err = workflow.HandleFlow(todos[0].RecordID, fdy, input)
	if err != nil {
		t.Fatal(err.Error())
	}

	// 查询待办
	todos, err = workflow.QueryTodoFlows(flowCode, bzr)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// 处理流程（通过）
	input["action"] = "pass"
	result, err = workflow.HandleFlow(todos[0].RecordID, bzr, input)
	if err != nil {
		t.Fatal(err.Error())
	}

	// 查询待办
	todos, err = workflow.QueryTodoFlows(flowCode, fdy)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// 处理流程（通过）
	input["action"] = "back"
	result, err = workflow.HandleFlow(todos[0].RecordID, fdy, input)
	if err != nil {
		t.Fatal(err.Error())
	}

	if result.NextNodes[0].CandidateIDs[0] != launcher {
		t.Fatalf("无效的下一级流转：%s", result.String())
	}
}

func TestQueryLastNodeInstance(t *testing.T) {
	result, err := workflow.QueryLastNodeInstance("b96558d1-d5e2-4cfe-8602-0dfd6b4be262")
	fmt.Printf("%v %v\n", result, err)
}


func TestForm(t *testing.T) {
	rst, err := workflow.StartFlow("Process_1", "StartEvent_1", "jmniu", nil)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	rst, err = workflow.HandleFlow(rst.FlowInstance.RecordID, "jmniu", nil)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
}

func TestProcForm(t *testing.T) {
	rst, err := workflow.StartFlow("proc_form", "id_start", "jmniu", map[string]interface{}{
		"name": "牛家明",
	})
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	rst, err = workflow.HandleFlow(rst.NextNodes[0].NodeInstance.RecordID, "jmniu", map[string]interface{}{
		"result": "同意",
	})
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
}