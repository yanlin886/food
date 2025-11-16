package order

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"

	"jihulab.com/yanlin/food-api/pkg/constants"
	"jihulab.com/yanlin/food-api/pkg/entity"
	"jihulab.com/yanlin/food-api/pkg/order/method"
	"jihulab.com/yanlin/food-api/pkg/test"
	"jihulab.com/yanlin/food-api/pkg/util"
)

const (
	TestCase1    = "新建产品订单 - 不含配料"
	TestCase1_1  = "新建产品订单 - 不含配料 - 检测入会门店"
	TestCase1_2  = "新建产品订单 - 不含配料 - 正餐 - 先吃后付 - 手动接单"
	TestCase1_3  = "新建产品订单 - 零支付配置"
	TestCase2    = "新建充值订单 - 含赠送金额，赠送积分，现金支付方式"
	TestCase2_1  = "取消充值订单"
	TestCase3    = "新建产品订单 - 使用积分抵扣"
	TestCase3_1  = "新建产品订单 - 使用积分抵扣_余额支付"
	TestCase4    = "新建产品订单 - 余额支付"
	TestCase5    = "新建产品订单 - 含配料"
	TestCase6    = "新建产品订单 - 支付方式为现金支付"
	TestCase7    = "新建充值订单 - 含赠送金额，赠送积分，现金支付方式, 会员升级"
	TestCase8    = "新建产品订单 - 含不存在的代金优惠券"
	TestCase9    = "新建产品订单 - 含存在的代金优惠券"
	TestCase10   = "新建产品订单 - 含存在的指定商品优惠券"
	TestCase11   = "小程序新建产品订单 - 不含配料, 前后端价格一致"
	TestCase11_1 = "小程序新建产品订单 - 不含配料, 前后端价格一致 - 余额支付成功"
	TestCase12   = "小程序新建产品订单 - 不含配料, 前后端价格不一致"
	TestCase13   = "获取产品订单 - 产品，管理端, 含优惠券，会员价"
	TestCase14   = "小程序订单列表"
	TestCase15   = "获取产品订单 - 产品，管理端, 含会员价"
	TestCase17   = "小程序订单详情"
	TestCase18   = "小程序充值 - 含赠送金额，赠送积分，现金支付方式, 会员升级"
	TestCase19   = "小程序订单详情 - 无订单流水号"
	TestCase20   = "web订单详情 - 无订单流水号"
	TestCase21   = "会员价下单测试"
	TestCase22   = "微信小程序储值明细列表"
	TestCase23   = "管理台储值记录"
	TestCase24   = "再来一单 历史订单查询"
	TestCase25   = "小程序新建产品订单 - 自提订单"
	TestCase26   = "再来一单 历史订单查询  配料为空"
	TestCase27   = "反结账"
	TestCase28   = "收银正餐下单测试"
	TestCase29   = "租户获取产品订单列表"
	TestCase30   = "租户获取储值订单列表"
	TestCase31   = "web端获取订单列表"
	TestCase32   = "新建产品订单 - 现金支付 - 校验订单消费记录"
	TestCase33   = "新建产品订单 - 余额支付 - 校验订单消费记录"
	TestCase34   = "新建产品订单 - 线上支付 - 校验订单消费记录"
	TestCase35   = "更新订单状态"

	TestCaseApp1 = "APP获取订单列表"
	TestCaseApp2 = "APP获取订单详情"
	TestCaseApp3 = "APP获取储值订单列表"
)

type VerifyTest struct {
	pool *pgxpool.Pool
}

func (v *VerifyTest) Switch(tc test.APITestCase, actualResponse string, t *testing.T) {
	switch tc.Name {
	case TestCase1:
		v.OrderProductWithoutIngredient(t, actualResponse)
	case TestCase1_1:
		v.TestCase1_1(t, actualResponse)
	case TestCase1_2:
		v.TestCase1_2(t, actualResponse)
	case TestCase2:
		v.OrderRechargeWithExtras(t, actualResponse, tc)
	case TestCase2_1:
		v.RechargeOrderCancel(t, actualResponse, tc)
	case TestCase3:
		v.OrderProductWithPoints(t, actualResponse)
	case TestCase3_1:
		v.OrderProductWithPointsAndBalance(t, actualResponse)
	case TestCase4:
		v.OrderProductWithBalance(t, actualResponse, tc)
	case TestCase5:
		v.OrderProductWithIngredient(t, actualResponse)
	case TestCase7:
		v.WebUpgradeMembership(t, tc)
	case TestCase9:
		v.Case9(t, tc)
		expectedCouponID := int64(1)
		v.OrderCoupon(t, actualResponse, tc, expectedCouponID)
	case TestCase10:
		v.Case10(t, actualResponse, tc)
		expectedCouponID := int64(2)
		v.OrderCoupon(t, actualResponse, tc, expectedCouponID)
	case TestCase18:
		v.WechatUpgradeMembership(t, tc)
	case TestCase21:
		v.MembershipPriceOrder(t, actualResponse, tc)
	case TestCase25:
		v.PickUpOrder(t, actualResponse, tc)
	case TestCase27:
		v.OrderReverse(t, actualResponse, tc)
	case TestCase28:
		v.TestCase28(t, actualResponse, tc)
	case TestCase32:
		v.TestCase32(t, actualResponse, tc)
	case TestCase33:
		v.TestCase33(t, actualResponse, tc)
	}
}

func (v *VerifyTest) OrderProductWithoutIngredient(t *testing.T, actualResponse string) {
	var (
		err   error
		query string
	)
	res := getCreateResponse{}
	err = json.Unmarshal([]byte(actualResponse), &res)
	require.Nil(t, err)
	orderID := res.Data.OrderID
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.TenantID, constants.MockTenantID)

	oip := struct {
		ProductID           int64
		SkuID               int64
		OriginalTotalAmount int64
	}{}

	// 验证order_item_products记录正确
	query = `select oip.product_id, oip.sku_id, oi.original_total_amount
		from order_item_products oip
		join order_items oi on oip.order_item_id = oi.id
		join orders o on o.id = oi.order_id
		where o.id = $1`
	err = pgxscan.Get(ctx, v.pool, &oip, query, orderID)
	require.Nil(t, err)
	require.Equal(t, int64(1), oip.ProductID)
	require.Equal(t, int64(2), oip.SkuID)
	require.Equal(t, int64(1800), oip.OriginalTotalAmount)

	// 验证支付记录
	pay := struct {
		Method string
		Amount int64
		Status string
	}{}
	query = `select method, amount, status
		from payments
		where order_id = $1`
	err = pgxscan.Get(ctx, v.pool, &pay, query, orderID)
	require.Nil(t, err, "test case OrderProductWithoutIngredient error getting payments")
	// @todo 支付方式是pending？
	require.Equal(t, "pending", pay.Method)
	require.Equal(t, int64(1800), pay.Amount)
	require.Equal(t, "pending", pay.Status)

	// 验证订单操作记录
	var operationInfoList []entity.OrderOperationForTest
	query = `select oo.id, ooa.admin_id,author_type, author_name, event, content
				from order_operations  oo
				left join order_operation_admins ooa on oo.id = ooa.order_operation_id
				where order_id = $1`
	err = pgxscan.Select(ctx, v.pool, &operationInfoList, query, orderID)
	require.NoError(t, err, t.Name())
	require.Equal(t, 1, len(operationInfoList), t.Name())

	expectedOperationInfoList := `[
			{
				"ID": 1,
				"AdminID": 1,
				"MembershipID": null,
				"AuthorType": "cashier",
				"AuthorName": "jack",
				"Event": "create",
				"Content": ""
			}
		]`
	jsonString, _ := json.Marshal(operationInfoList)

	require.JSONEq(t, expectedOperationInfoList, string(jsonString))

	// 验证membership 1的入会门店id是1.
	var (
		entryStoreID   int64
		entryStoreName string
	)
	query = `select entry_store_id, entry_store_name from memberships where id = $1`
	err = v.pool.QueryRow(ctx, query, 1).Scan(&entryStoreID, &entryStoreName)
	require.Nil(t, err)
	require.Equal(t, int64(1), entryStoreID)
	require.Equal(t, "金鹰门店1", entryStoreName)

	// 验证快餐的桌台号
	var tableID string
	query = `select table_id from order_table_ids where order_id = $1`
	err = v.pool.QueryRow(ctx, query, orderID).Scan(&tableID)
	require.Nil(t, err, "Test case 1, 桌台号")
	require.Equal(t, "A1", tableID)
}

func (v *VerifyTest) TestCase1_1(t *testing.T, actualResponse string) {
	var (
		err   error
		query string
	)
	// 验证membership 1的入会门店id是4.
	var (
		entryStoreID   int64
		entryStoreName string
	)
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.TenantID, constants.MockTenantID)
	query = `select entry_store_id, entry_store_name from memberships where id = $1`
	err = v.pool.QueryRow(ctx, query, 1).Scan(&entryStoreID, &entryStoreName)
	require.Nil(t, err)
	require.Equal(t, int64(4), entryStoreID)
	require.Equal(t, "金鹰门店4", entryStoreName)
}

func (v *VerifyTest) TestCase1_2(t *testing.T, actualResponse string) {
	var (
		err   error
		query string
	)
	res := getCreateResponse{}
	err = json.Unmarshal([]byte(actualResponse), &res)
	require.Nil(t, err)
	orderID := res.Data.OrderID
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.TenantID, constants.MockTenantID)

	oip := struct {
		TakeStatus string
	}{}

	query = `select take_status
		from order_item_products oip
		join order_items oi on oip.order_item_id = oi.id
		join orders o on o.id = oi.order_id
		where o.id = $1`
	err = pgxscan.Get(ctx, v.pool, &oip, query, orderID)
	require.Nil(t, err)
	require.Equal(t, "pending_take_order", oip.TakeStatus)
}

func (v *VerifyTest) OrderRechargeWithExtras(t *testing.T, actualResponse string, tc test.APITestCase) {
	by := struct {
		MembershipID  string `json:"membership_id"`
		RechargeID    string `json:"recharge_id"`
		RechargeStall struct {
			ID             string             `json:"id"`
			RewardAmount   entity.MoneyString `json:"reward_amount"`
			RewardPoints   int64              `json:"reward_points"`
			RechargeAmount entity.MoneyString `json:"recharge_amount"`
		} `json:"recharge_stall"`
	}{}

	res := getCreateResponse{}
	err := json.Unmarshal([]byte(actualResponse), &res)
	require.Nil(t, err, "test case OrderRechargeWithExtras error unmarshalling response")

	err = json.Unmarshal([]byte(tc.Body), &by)
	require.Nil(t, err, "test case OrderRechargeWithExtras error unmarshalling request")

	orderID := res.Data.OrderID
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.TenantID, constants.MockTenantID)

	// 验证是否创建了membership_balance_recharges记录
	query := `select cast(membership_id as text), beneficial_store_id, principal_amount, reward_amount
		from membership_balance_recharges
		where recharge_order_id = $1`
	mbr := struct {
		MembershipID      string
		BeneficialStoreID int64
		PrincipalAmount   entity.MoneyInt64
		RewardAmount      entity.MoneyInt64
	}{}
	err = pgxscan.Get(ctx, v.pool, &mbr, query, orderID)
	require.Nil(t, err, "test case OrderRechargeWithExtras error getting membership_balance_recharges")
	require.Equal(t, by.MembershipID, mbr.MembershipID)
	require.Equal(t, int64(1), mbr.BeneficialStoreID)
	require.Equal(t, by.RechargeStall.RechargeAmount, mbr.PrincipalAmount.MustMoneyString())
	require.Equal(t, by.RechargeStall.RewardAmount, mbr.RewardAmount.MustMoneyString())

	oir := struct {
		CurrentBalanceAmount         int64
		CurrentPrincipalAmount       int64
		CurrentRewardAmount          int64
		TotalRechargeAmount          int64
		TotalRechargePrincipalAmount int64
		TotalRewardAmount            int64
	}{}
	// 验证order_item_recharges记录正确
	query = `select 
    		current_balance_amount,
       		current_principal_amount,
       		current_reward_amount,
       		total_recharge_amount,
       		total_recharge_principal_amount,
       		total_reward_amount
       from membership_balances where membership_id = $1`
	err = pgxscan.Get(ctx, v.pool, &oir, query, by.MembershipID)
	require.Nil(t, err, "test case OrderRechargeWithExtras error getting membership_balances")

	rechargeAmount, err := by.RechargeStall.RechargeAmount.Int64()
	require.Nil(t, err, "test case OrderRechargeWithExtras error getting recharge amount")
	rewardAmount, err := by.RechargeStall.RewardAmount.Int64()
	require.Nil(t, err, "test case OrderRechargeWithExtras error getting reward amount")
	// 验证生成数据
	require.Equal(t, rechargeAmount+rewardAmount+200, oir.CurrentBalanceAmount)
	require.Equal(t, int64(7400), oir.CurrentBalanceAmount)
	require.Equal(t, rechargeAmount+100, oir.CurrentPrincipalAmount)
	require.Equal(t, rewardAmount+100, oir.CurrentRewardAmount)
	require.Equal(t, rechargeAmount+rewardAmount+200, oir.TotalRechargeAmount)
	require.Equal(t, rechargeAmount+100, oir.TotalRechargePrincipalAmount)
	require.Equal(t, rewardAmount+100, oir.TotalRewardAmount)

	// 验证积分
	query = `select points from membership_points where membership_id = $1`
	data := struct {
		Points int64
	}{}
	err = pgxscan.Get(ctx, v.pool, &data, query, by.MembershipID)
	require.Nil(t, err)
	require.Equal(t, by.RechargeStall.RewardPoints+100, data.Points)

	// 验证订单操作记录
	var operationInfoList []entity.OrderOperationForTest
	query = `select oo.id, ooa.admin_id,author_type, author_name, event, content
				from order_operations  oo
				left join order_operation_admins ooa on oo.id = ooa.order_operation_id
				where order_id = $1`
	err = pgxscan.Select(ctx, v.pool, &operationInfoList, query, orderID)
	require.NoError(t, err, t.Name())
	require.Equal(t, 2, len(operationInfoList), t.Name())

	expectedOperationInfoList := `[
			{
				"ID": 1,
				"AdminID": 1,
				"MembershipID": null,
				"AuthorType": "cashier",
				"AuthorName": "jack",
				"Event": "create",
				"Content": ""
			},
			{
				"ID": 2,
				"AdminID": 1,
				"MembershipID": null,
				"AuthorType": "cashier",
				"AuthorName": "jack",
				"Event": "pay_success",
				"Content": "收银端充值：充值金额￥31"
			}
		]`
	jsonString, _ := json.Marshal(operationInfoList)

	require.JSONEq(t, expectedOperationInfoList, string(jsonString), tc.Name)

	// 确保订单的total_amount跟actual_amount相同
	var totalAmount, actualAmount int64
	query = `select total_amount, actual_amount from orders where id = $1`
	err = v.pool.QueryRow(ctx, query, orderID).Scan(&totalAmount, &actualAmount)
	require.Nil(t, err)
	require.Equal(t, int64(3100), totalAmount, "order total amount")
	require.Equal(t, int64(3100), actualAmount, "order actual amount")

	// 验证order_memberships的当前余额是支付后的余额
	var currentBalance int64
	query = `select current_balance_amount from order_memberships where order_id = $1`
	err = v.pool.QueryRow(ctx, query, orderID).Scan(&currentBalance)
	require.Nil(t, err)
	require.Equal(t, int64(7400), currentBalance, "order_memberships的当前余额")

	// 记录order_recharges的赠送金额
	var rechargeRewardAmount int64
	query = `select reward_amount from order_recharges where order_id = $1`
	err = v.pool.QueryRow(ctx, query, orderID).Scan(&rechargeRewardAmount)
	require.Nil(t, err)
	require.Equal(t, int64(4100), rechargeRewardAmount)
}

func (v *VerifyTest) OrderProductWithBalance(t *testing.T, actualResponse string, tc test.APITestCase) {
	by := struct {
		MembershipID string `json:"membership_id"`
		RechargeID   string `json:"recharge_id"`
	}{}

	err := json.Unmarshal([]byte(tc.Body), &by)
	require.Nil(t, err)

	res := struct {
		Data struct {
			OrderID string `json:"order_id"`
			Status  string `json:"status"`
			Type    string `json:"type"`
		} `json:"data"`
	}{}

	err = json.Unmarshal([]byte(actualResponse), &res)
	require.Nil(t, err)

	// 获取用户余额数据
	query := `select 
    		  current_balance_amount,
       		  current_principal_amount,
       		  current_reward_amount
		   	  from membership_balances where membership_id = $1`

	data := struct {
		CurrentBalanceAmount   int64 `json:"current_balance_amount"`
		CurrentPrincipalAmount int64 `json:"current_principal_amount"`
		CurrentRewardAmount    int64 `json:"current_reward_amount"`
	}{}

	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.TenantID, constants.MockTenantID)
	err = pgxscan.Get(ctx, v.pool, &data, query, by.MembershipID)
	require.Nil(t, err)

	// 校验余额数据
	require.Equal(t, int64(100), data.CurrentPrincipalAmount)
	require.Equal(t, int64(97), data.CurrentRewardAmount)
	require.Equal(t, int64(197), data.CurrentBalanceAmount)

	query = `select total_consume_amount,consume_count from membership_consumes where membership_id = $1`
	consume := struct {
		TotalConsumeAmount int64 `json:"total_consume_amount"`
		ConsumeCount       int64 `json:"consume_count"`
	}{}

	// 校验消费记录数据
	err = pgxscan.Get(ctx, v.pool, &consume, query, by.MembershipID)
	require.Nil(t, err)
	require.Equal(t, int64(103), consume.TotalConsumeAmount)
	require.Equal(t, int64(2), consume.ConsumeCount)

	// 校验订单状态
	var status string
	query = `select status from orders where id = $1`
	err = pgxscan.Get(ctx, v.pool, &status, query, res.Data.OrderID)
	require.Nil(t, err)
	require.Equal(t, "completed", status)

	// 验证支付状态
	var paymentStatus string
	query = `select status from payments where order_id = $1`
	err = pgxscan.Get(ctx, v.pool, &paymentStatus, query, res.Data.OrderID)
	require.Nil(t, err)
	require.Equal(t, "paid", paymentStatus)

	// 验证订单操作记录
	var operationInfoList []entity.OrderOperationForTest
	query = `select oo.id, ooa.admin_id,author_type, author_name, event, content
				from order_operations  oo
				left join order_operation_admins ooa on oo.id = ooa.order_operation_id
				where order_id = $1`
	err = pgxscan.Select(ctx, v.pool, &operationInfoList, query, res.Data.OrderID)
	require.NoError(t, err, t.Name())
	require.Equal(t, 2, len(operationInfoList), t.Name())

	expectedOperationInfoList := `[
			{
				"ID": 1,
				"AdminID": 1,
				"MembershipID": null,
				"AuthorType": "cashier",
				"AuthorName": "jack",
				"Event": "create",
				"Content": ""
			},
			{
				"ID": 2,
				"AdminID": 1,
				"MembershipID": null,
				"AuthorType": "cashier",
				"AuthorName": "jack",
				"Event": "pay_success",
				"Content": "结账：订单金额￥0.03, 优惠合计￥0, 实收金额￥0.03(余额￥0.03)"
			}
		]`
	jsonString, _ := json.Marshal(operationInfoList)
	require.JSONEq(t, expectedOperationInfoList, string(jsonString), tc.Name)
}

func (v *VerifyTest) OrderProductWithPoints(t *testing.T, actualResponse string) {
	res := getCreateResponse{}
	err := json.Unmarshal([]byte(actualResponse), &res)
	require.Nil(t, err)
	orderID := res.Data.OrderID
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.TenantID, constants.MockTenantID)

	// 验证是否有积分扣减记录
	decrease := struct {
		MembershipID  int64
		Points        int64
		CurrentPoints int64
		CashOut       int64
	}{}
	query := `select membership_id, points, current_points, cash_out
		from order_point_decreases
		where order_id = $1`
	err = pgxscan.Get(ctx, v.pool, &decrease, query, orderID)
	require.NoError(t, err)
	require.Equal(t, int64(10), decrease.Points)
	require.Equal(t, int64(90), decrease.CurrentPoints)
	require.Equal(t, int64(200), decrease.CashOut)

	// 验证会员积分是否发生正确的变化
	query = `select points
	from membership_points
	where membership_id = $1`
	var points int64
	err = pgxscan.Get(ctx, v.pool, &points, query, decrease.MembershipID)
	require.NoError(t, err)
	// 会员增加的积分只有在支付成功后才会增加
	// 原先是100积分，变化后应该是100-10=90
	require.Equal(t, int64(90), points)
	// 验证订单操作记录
	var operationInfoList []entity.OrderOperationForTest
	query = `select oo.id, ooa.admin_id,author_type, author_name, event, content
				from order_operations  oo
				left join order_operation_admins ooa on oo.id = ooa.order_operation_id
				where order_id = $1`
	err = pgxscan.Select(ctx, v.pool, &operationInfoList, query, res.Data.OrderID)
	require.NoError(t, err, t.Name())
	require.Equal(t, 1, len(operationInfoList), t.Name())

	expectedOperationInfoList := `[
			{
				"ID": 1,
				"AdminID": 1,
				"MembershipID": null,
				"AuthorType": "cashier",
				"AuthorName": "jack",
				"Event": "create",
				"Content": ""
			}
		]`
	jsonString, _ := json.Marshal(operationInfoList)
	require.JSONEq(t, expectedOperationInfoList, string(jsonString))
}

func (v *VerifyTest) OrderProductWithPointsAndBalance(t *testing.T, actualResponse string) {
	res := getCreateResponse{}
	err := json.Unmarshal([]byte(actualResponse), &res)
	require.Nil(t, err)
	orderID := res.Data.OrderID
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.TenantID, constants.MockTenantID)

	// 验证是否有积分扣减记录
	decrease := struct {
		MembershipID  int64
		Points        int64
		CurrentPoints int64
		CashOut       int64
	}{}
	query := `select membership_id, points, current_points, cash_out
		from order_point_decreases
		where order_id = $1`
	err = pgxscan.Get(ctx, v.pool, &decrease, query, orderID)
	require.NoError(t, err)
	require.Equal(t, int64(10), decrease.Points)
	require.Equal(t, int64(90), decrease.CurrentPoints)
	require.Equal(t, int64(200), decrease.CashOut)
	// 验证是否有积分增加记录，且增加的积分是正确的
	increase := struct {
		Points        int64
		CurrentPoints int64
	}{}
	query = `select points, current_points
		from order_point_increases
		where order_id = $1`
	err = pgxscan.Get(ctx, v.pool, &increase, query, orderID)
	require.NoError(t, err)
	require.Equal(t, int64(4), increase.Points)
	require.Equal(t, int64(94), increase.CurrentPoints)
	// 验证会员积分是否发生正确的变化
	query = `select points
	from membership_points
	where membership_id = $1`
	var points int64
	err = pgxscan.Get(ctx, v.pool, &points, query, decrease.MembershipID)
	require.NoError(t, err)
	// 会员增加的积分只有在支付成功后才会增加
	// 原先是100积分，变化后应该是100-10+4=94
	require.Equal(t, int64(94), points)

	// 验证订单操作记录
	var operationInfoList []entity.OrderOperationForTest
	query = `select oo.id, ooa.admin_id,author_type, author_name, event, content
				from order_operations  oo
				left join order_operation_admins ooa on oo.id = ooa.order_operation_id
				where order_id = $1`
	err = pgxscan.Select(ctx, v.pool, &operationInfoList, query, res.Data.OrderID)
	require.NoError(t, err, t.Name())
	require.Equal(t, 2, len(operationInfoList), t.Name())

	expectedOperationInfoList := `[
			{
				"ID": 1,
				"AdminID": 1,
				"MembershipID": null,
				"AuthorType": "cashier",
				"AuthorName": "jack",
				"Event": "create",
				"Content": ""
			},
			{
				"ID": 2,
				"AdminID": 1,
				"MembershipID": null,
				"AuthorType": "cashier",
				"AuthorName": "jack",
				"Event": "pay_success",
				"Content": "结账：订单金额￥4, 优惠合计￥2, 实收金额￥2(余额￥2)"
			}
		]`
	jsonString, _ := json.Marshal(operationInfoList)
	require.JSONEq(t, expectedOperationInfoList, string(jsonString))
}

func (v *VerifyTest) OrderProductWithIngredient(t *testing.T, actualResponse string) {
	var (
		err   error
		query string
	)
	res := getCreateResponse{}
	err = json.Unmarshal([]byte(actualResponse), &res)
	require.Nil(t, err)
	orderID := res.Data.OrderID
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.TenantID, constants.MockTenantID)

	oip := struct {
		OrderItemID int64
		ProductID   int64
		SkuID       int64
	}{}

	// 验证order_item_products记录正确
	query = `select oip.product_id, oip.sku_id, oi.id as order_item_id
		from order_item_products oip
		join order_items oi on oip.order_item_id = oi.id
		join orders o on o.id = oi.order_id
		where o.id = $1`
	err = pgxscan.Get(ctx, v.pool, &oip, query, orderID)
	require.Nil(t, err)
	require.Equal(t, int64(1), oip.ProductID)
	require.Equal(t, int64(2), oip.SkuID)

	// 验证order_item_ingredients记录正确
	ing := struct {
		OrderItemID        int64
		IngredientID       int64
		IngredientOptionID int64
	}{}

	query = `select oii.ingredient_id, oii.ingredient_option_id, oi.id as order_item_id
		from order_item_ingredients oii
		join order_items oi on oii.order_item_id = oi.id
		join orders o on o.id = oi.order_id
		where o.id = $1`
	err = pgxscan.Get(ctx, v.pool, &ing, query, orderID)
	require.Nil(t, err)
	require.Equal(t, int64(1), ing.IngredientID)
	require.Equal(t, int64(1), ing.IngredientOptionID)

	// 验证order_item_self_relations关联记录是否存在
	var exists bool
	query = `select exists
    (select 1 from order_item_self_relations 
              where order_item_id = $1 and order_item_parent_id = $2)`
	err = pgxscan.Get(ctx, v.pool, &exists, query, ing.OrderItemID, oip.OrderItemID)
	require.Nil(t, err)
	require.True(t, exists)

	// 验证支付记录
	pay := struct {
		Method string
		Amount int64
		Status string
	}{}
	query = `select method, amount, status
		from payments
		where order_id = $1`
	err = pgxscan.Get(ctx, v.pool, &pay, query, orderID)
	require.Nil(t, err, "test case OrderProductWithIngredient error getting payments")
	require.Equal(t, "pending", pay.Method)
	require.Equal(t, int64(2000), pay.Amount)
	require.Equal(t, "pending", pay.Status)

	// 验证订单操作记录
	var operationInfoList []entity.OrderOperationForTest
	query = `select oo.id, ooa.admin_id,author_type, author_name, event, content
				from order_operations  oo
				left join order_operation_admins ooa on oo.id = ooa.order_operation_id
				where order_id = $1`
	err = pgxscan.Select(ctx, v.pool, &operationInfoList, query, res.Data.OrderID)
	require.NoError(t, err, t.Name())
	require.Equal(t, 1, len(operationInfoList), t.Name())

	expectedOperationInfoList := `[
			{
				"ID": 1,
				"AdminID": 1,
				"MembershipID": null,
				"AuthorType": "cashier",
				"AuthorName": "jack",
				"Event": "create",
				"Content": ""
			}
		]`
	jsonString, _ := json.Marshal(operationInfoList)
	require.JSONEq(t, expectedOperationInfoList, string(jsonString))
}

func (v *VerifyTest) WebUpgradeMembership(t *testing.T, tc test.APITestCase) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.TenantID, constants.MockTenantID)
	query := `SELECT mcl."level" 
				FROM membership_membership_card_level_relations mmclr
				LEFT JOIN membership_card_levels mcl ON mcl.ID = mmclr.membership_card_level_id
				WHERE mmclr.membership_id = $1`

	var level int64 = -1
	by := struct {
		MembershipID string `json:"membership_id"`
	}{}

	err := json.Unmarshal([]byte(tc.Body), &by)
	require.Nil(t, err, "test case UpgradeMembership error unmarshalling body")

	err = pgxscan.Get(ctx, v.pool, &level, query, by.MembershipID)
	require.Nil(t, err)
	require.Equal(t, int64(10), level)
}

func (v *VerifyTest) WechatUpgradeMembership(t *testing.T, tc test.APITestCase) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.TenantID, constants.MockTenantID)
	query := `SELECT mcl."level" 
				FROM membership_membership_card_level_relations mmclr
				LEFT JOIN membership_card_levels mcl ON mcl.ID = mmclr.membership_card_level_id
				WHERE mmclr.membership_id = $1`

	var level int64 = -1

	err := pgxscan.Get(ctx, v.pool, &level, query, constants.MockMembershipID)
	require.Nil(t, err)
	require.Equal(t, int64(10), level)
}

func (v *VerifyTest) OrderCoupon(t *testing.T, actualResponse string, tc test.APITestCase, expectCouponID int64) {
	var (
		err   error
		query string
		oc    struct {
			MembershipID int64
			CouponSendID int64
			CouponID     int64
			ReduceAmount int64
			SaleAmount   *int64
		}
	)
	res := getCreateResponse{}
	err = json.Unmarshal([]byte(actualResponse), &res)
	require.Nil(t, err)
	orderID := res.Data.OrderID
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.TenantID, constants.MockTenantID)

	// order coupons
	query = `select membership_id, coupon_send_id, coupon_id, reduce_amount, sale_amount
		from order_coupons
		where order_id = $1`
	err = pgxscan.Get(ctx, v.pool, &oc, query, orderID)
	require.Nil(t, err)

	require.Equal(t, int64(1), oc.MembershipID)
	require.Equal(t, int64(1), oc.CouponSendID)
	require.Equal(t, expectCouponID, oc.CouponID)
	require.Equal(t, int64(1278), oc.ReduceAmount)
	require.Equal(t, int64(1800), *oc.SaleAmount)
}

func (v *VerifyTest) Case9(t *testing.T, tc test.APITestCase) {
	var (
		err    error
		query  string
		number int64
	)
	body := method.ProductRequest{}
	err = json.Unmarshal([]byte(tc.Body), &body)
	require.Nil(t, err)

	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.TenantID, constants.MockTenantID)

	// membership coupons
	membershipCouponID := body.CouponProperty.MembershipCouponID
	query = `select number from membership_coupons where id = $1`
	err = v.pool.QueryRow(ctx, query, membershipCouponID).Scan(&number)
	require.Nil(t, err)

	require.Equal(t, int64(1), number)
}

func (v *VerifyTest) Case10(t *testing.T, actualResponse string, tc test.APITestCase) {
	var (
		err    error
		query  string
		count  int64
		exists bool

		oics []struct {
			OriginalTotalAmount int64
			CouponName          string
			ReduceAmount        int64
		}
	)
	body := method.ProductRequest{}
	err = json.Unmarshal([]byte(tc.Body), &body)
	require.Nil(t, err)

	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.TenantID, constants.MockTenantID)

	// membership coupons
	membershipCouponID := body.CouponProperty.MembershipCouponID
	query = `select exists (select 1 from membership_coupons where id = $1)`
	err = v.pool.QueryRow(ctx, query, membershipCouponID).Scan(&exists)
	require.Nil(t, err)

	require.False(t, exists)

	res := getCreateResponse{}
	err = json.Unmarshal([]byte(actualResponse), &res)
	require.Nil(t, err)
	query = `select count(*) from order_item_coupons`
	err = v.pool.QueryRow(ctx, query).Scan(&count)
	require.Nil(t, err)
	require.Equal(t, int64(1), count, "数量有误1")

	// order item coupon
	query = `select oic.coupon_name, oic.reduce_amount, oi.original_total_amount
		from order_items oi
		left join order_item_coupons oic on oi.id = oic.order_item_id
		where oi.order_id = $1`
	err = pgxscan.Select(ctx, v.pool, &oics, query, res.Data.OrderID)
	require.Nil(t, err)

	require.Equal(t, 1, len(oics), "数量有误2")
	require.Equal(t, body.CouponProperty.Name, oics[0].CouponName)
	require.Equal(t, body.CouponProperty.ReduceAmount.MustInt64(), oics[0].ReduceAmount, "价格有误")
	require.Equal(t, int64(2000), oics[0].OriginalTotalAmount)
}

func (v *VerifyTest) MembershipPriceOrder(t *testing.T, actualResponse string, tc test.APITestCase) {
	query := `select 
				actual_amount,
				total_amount,
				id 
				from orders`
	ctx, cancelFunc := test.NewContext()
	defer cancelFunc()
	var actualAmount, totalAmount, orderID, reduceAmount int64
	err := v.pool.QueryRow(ctx, query).Scan(&actualAmount, &totalAmount, &orderID)
	require.Nil(t, err)
	require.Equal(t, int64(1800), totalAmount)
	require.Equal(t, int64(1200), actualAmount)

	query = `select reduce_amount from order_membership_reduce_amounts where order_id = $1`
	err = v.pool.QueryRow(ctx, query, orderID).Scan(&reduceAmount)
	require.Nil(t, err)
	require.Equal(t, int64(600), reduceAmount)

}
func (v *VerifyTest) OrderReverse(t *testing.T, actualResponse string, tc test.APITestCase) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.TenantID, constants.MockTenantID)

	orderID := 1670710200224452608

	// 查询订单的会员
	query := `select membership_id from order_memberships where order_id = $1`
	orderMembership := &entity.OrderMembership{}

	err := pgxscan.Get(ctx, v.pool, orderMembership, query, orderID)
	require.Nil(t, err)

	// 查询会员余额数据
	query = `select 
			current_balance_amount,
			current_principal_amount,
			current_reward_amount 
			from membership_balances 
			where membership_id = $1`
	recharge := struct {
		CurrentBalanceAmount   int64 `db:"current_balance_amount"`
		CurrentPrincipalAmount int64 `db:"current_principal_amount"`
		CurrentRewardAmount    int64 `db:"current_reward_amount"`
	}{}
	err = pgxscan.Get(ctx, v.pool, &recharge, query, orderMembership.MembershipID)
	require.Nil(t, err)
	require.Equal(t, int64(203), recharge.CurrentBalanceAmount)
	require.Equal(t, int64(100), recharge.CurrentPrincipalAmount)
	require.Equal(t, int64(103), recharge.CurrentRewardAmount)

	query = `select 
            total_consume_amount,
			consume_count
			from membership_consumes
			where membership_id = $1`
	consume := struct {
		TotalConsumeAmount int64 `db:"total_consume_amount"`
		ConsumeCount       int64 `db:"consume_count"`
	}{}
	err = pgxscan.Get(ctx, v.pool, &consume, query, orderMembership.MembershipID)
	require.Nil(t, err)
	require.Equal(t, int64(100), consume.TotalConsumeAmount)
	require.Equal(t, int64(1), consume.ConsumeCount)

	// 验证订单操作记录
	var operationInfoList []entity.OrderOperationForTest
	query = `select oo.id, ooa.admin_id,author_type, author_name, event, content
				from order_operations  oo
				left join order_operation_admins ooa on oo.id = ooa.order_operation_id
				where order_id = $1`
	err = pgxscan.Select(ctx, v.pool, &operationInfoList, query, orderID)
	require.NoError(t, err, t.Name())
	require.Equal(t, 1, len(operationInfoList), t.Name())

	expectedOperationInfoList := `[
			{
				"ID": 1,
				"AdminID": 1,
				"MembershipID": null,
				"AuthorType": "cashier",
				"AuthorName": "jack",
				"Event": "reverse",
				"Content": "反结账：退款金额：0.03(余额￥0.03) <br/> 退款原因：点错了"
			}
		]`
	jsonString, _ := json.Marshal(operationInfoList)

	require.JSONEq(t, expectedOperationInfoList, string(jsonString), tc.Name)
}

func (v *VerifyTest) PickUpOrder(t *testing.T, actualResponse string, tc test.APITestCase) {
	// 这里主要检查一下自提订单的自提时间是否正确
	var (
		query               string
		actualDiningTime    entity.JSONTime
		expectDiningTimeStr = entity.TimeString("2023-09-15 17:00:00")
	)

	expectDiningTime, err := expectDiningTimeStr.UTCJSONTime()
	require.Nil(t, err)

	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.TenantID, constants.MockTenantID)

	res := getCreateResponse{}
	err = json.Unmarshal([]byte(actualResponse), &res)
	require.Nil(t, err)

	query = `select dining_time from orders where id = $1`
	err = v.pool.QueryRow(ctx, query, res.Data.OrderID).Scan(&actualDiningTime)
	require.Nil(t, err)
	require.Equal(t, expectDiningTime, actualDiningTime)

	// 验证订单操作记录
	var operationInfoList []entity.OrderOperationForTest
	query = `select oo.id, oom.membership_id,author_type, author_name, event, content
				from order_operations  oo
				left join order_operation_memberships oom on oo.id = oom.order_operation_id
				where order_id = $1`
	err = pgxscan.Select(ctx, v.pool, &operationInfoList, query, res.Data.OrderID)
	require.NoError(t, err, tc.Name)
	require.Equal(t, 3, len(operationInfoList), tc.Name)

	expectedOperationInfoList := `[
			{
				"ID": 1,
				"AdminID": null,
				"MembershipID": 1,
				"AuthorType": "membership",
				"AuthorName": "小林",
				"Event": "create",
				"Content": ""
			},
			{
				"ID": 2,
				"AdminID": null,
				"MembershipID": 1,
				"AuthorType": "membership",
				"AuthorName": "小林",
				"Event": "pay_success",
				"Content": "支付成功：订单金额￥18, 优惠合计￥0, 实收金额￥18(现金￥18)"
			},
			{
				"ID": 3,
				"AdminID": null,
				"MembershipID": null,
				"AuthorType": "system",
				"AuthorName": "system",
				"Event": "receiving_order",
				"Content": "接单"
			}
		]`
	jsonString, _ := json.Marshal(operationInfoList)

	require.JSONEq(t, expectedOperationInfoList, string(jsonString))
}

func (v *VerifyTest) RechargeOrderCancel(t *testing.T, actualResponse string, tc test.APITestCase) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.TenantID, constants.MockTenantID)

	// 验证订单操作记录
	var operationInfoList []entity.OrderOperationForTest
	query := `select oo.id, ooa.admin_id,author_type, author_name, event, content
				from order_operations  oo
				left join order_operation_admins ooa on oo.id = ooa.order_operation_id
				where order_id = $1`
	err := pgxscan.Select(ctx, v.pool, &operationInfoList, query, 1670710200224452612)
	require.NoError(t, err, tc.Name)
	require.Equal(t, 1, len(operationInfoList), tc.Name)

	expectedOperationInfoList := `[
			{
				"ID": 1,
				"AdminID": 1,
				"MembershipID": null,
				"AuthorType": "cashier",
				"AuthorName": "jack",
				"Event": "cancel",
				"Content": "取消订单 <br/> 原因：不想要了"
			}
		]`
	jsonString, _ := json.Marshal(operationInfoList)

	require.JSONEq(t, expectedOperationInfoList, string(jsonString), tc.Name)
}

func (v *VerifyTest) TestCase28(t *testing.T, actualResponse string, tc test.APITestCase) {
	query := `select table_id,name from order_tables`
	ctx, cancelFunc := test.NewContext()
	defer cancelFunc()
	var (
		serveiceFeeID int64
		amount        int
		tableID       int64
		name          string
	)
	err := v.pool.QueryRow(ctx, query).Scan(&tableID, &name)
	require.Nil(t, err)
	require.Equal(t, int64(1), tableID)
	require.Equal(t, "台1", name)

	query = `select service_fee_id,name,amount from order_service_fees`
	err = v.pool.QueryRow(ctx, query).Scan(&serveiceFeeID, &name, &amount)
	require.Nil(t, err)
	require.Equal(t, int64(1), serveiceFeeID)
	require.Equal(t, "服务费1", name)
	require.Equal(t, 10010, amount)

}

func (v *VerifyTest) TestCase32(t *testing.T, actualResponse string, tc test.APITestCase) {
	var (
		err   error
		query string
		body  method.ProductRequest
		res   getCreateResponse
	)

	err = json.Unmarshal([]byte(actualResponse), &res)
	require.Nil(t, err, "error unmarshalling response")

	err = json.Unmarshal([]byte(tc.Body), &body)
	require.Nil(t, err, "error unmarshalling body")

	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.TenantID, constants.MockTenantID)

	orderConsumes := entity.OrderConsume{}
	consumeOrderID, err := util.ConvertStringInt64(res.Data.OrderID)
	require.Nil(t, err, "error converting string to int64")

	query = `select * from order_consumes where consume_order_id = $1`
	err = pgxscan.Get(ctx, v.pool, &orderConsumes, query, consumeOrderID)
	require.Nil(t, err, "error getting order consumes")

	require.Equal(t, constants.PaymentMethodCash, orderConsumes.PayMethod)
	require.Equal(t, body.ActualAmount.MustMoneyInt64(), orderConsumes.ConsumeAmount)
	require.Equal(t, body.ActualAmount.MustMoneyInt64(), orderConsumes.PrincipalAmount)
	require.Equal(t, entity.MoneyInt64(0), orderConsumes.RewardAmount)
	require.Equal(t, consumeOrderID, orderConsumes.ConsumeOrderID)
	require.Equal(t, int64(1), orderConsumes.ConsumeStoreID)
	require.Nil(t, orderConsumes.RechargeOrderID)
	require.Nil(t, orderConsumes.MembershipID)
}

func (v *VerifyTest) TestCase33(t *testing.T, actualResponse string, tc test.APITestCase) {
	var (
		query            string
		balanceRecharges []struct {
			ID                int64 `json:"id"`
			PrincipalAmount   int64 `json:"principal_amount"`
			RewardAmount      int64 `json:"reward_amount"`
			BeneficialStoreID int64 `json:"beneficial_store_id"`
			RechargeOrderID   int64 `json:"recharge_order_id"`
		}
		// res getCreateResponse
		res struct {
			Data struct {
				OrderID string `json:"order_id"`
			} `json:"data"`
		}
	)

	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.TenantID, constants.MockTenantID)
	query = `select id, principal_amount, reward_amount, beneficial_store_id, recharge_order_id
		from membership_balance_recharges where membership_id = $1 order by id ASC`
	membershipID := 1
	err := pgxscan.Select(ctx, v.pool, &balanceRecharges, query, membershipID)
	require.Nil(t, err)
	jsonRes, err := json.Marshal(balanceRecharges)
	require.Nil(t, err)
	// 需要扣除 principal(40), reward(40)
	expectedJson := `
	[
		{
			"id": 2,
			"principal_amount": 20,
			"reward_amount": 0,
			"beneficial_store_id": 4,
			"recharge_order_id": 1670710200224452609
		},
		{
			"id": 3,
			"principal_amount": 40,
			"reward_amount": 60,
			"beneficial_store_id": 5,
			"recharge_order_id": 1670710200224452610
		}
	]`
	require.JSONEq(t, expectedJson, string(jsonRes))

	err = json.Unmarshal([]byte(actualResponse), &res)
	require.Nil(t, err, "error unmarshalling response")

	orderConsumes := []struct {
		ID                int64                   `json:"id"`
		TenantID          string                  `json:"-"`
		PayMethod         constants.PaymentMethod `json:"pay_method"`
		PrincipalAmount   int64                   `json:"principal_amount"`
		RewardAmount      int64                   `json:"reward_amount"`
		ConsumeAmount     int64                   `json:"consume_amount"`
		ConsumeOrderID    int64                   `json:"-"`
		ConsumeStoreID    int64                   `json:"consume_store_id"`
		BeneficialStoreID int64                   `json:"beneficial_store_id"`
		RechargeOrderID   *int64                  `json:"recharge_order_id"`
		MembershipID      *int64                  `json:"membership_id"`
		CreatedAt         entity.JSONTime         `json:"-"`
	}{}
	consumeOrderID, err := util.ConvertStringInt64(res.Data.OrderID)
	require.Nil(t, err, "error converting string to int64")

	query = `select * from order_consumes where consume_order_id = $1`
	err = pgxscan.Select(ctx, v.pool, &orderConsumes, query, consumeOrderID)
	require.Nil(t, err, "error getting order consumes")

	jsonRes, err = json.Marshal(orderConsumes)
	require.Nil(t, err)
	expectedJson = `[
		{
			"id": 1,
			"pay_method": "balance",
			"principal_amount": 30,
			"reward_amount": 10,
			"consume_amount": 80,
			"consume_store_id": 1,
			"beneficial_store_id": 1,
			"recharge_order_id": 1670710200224452608,
			"membership_id": 1
		},
		{
			"id": 2,
			"pay_method": "balance",
			"principal_amount": 10,
			"reward_amount": 10,
			"consume_amount": 80,
			"consume_store_id": 1,
			"beneficial_store_id": 4,
			"recharge_order_id": 1670710200224452609,
			"membership_id": 1
		},
		{
			"id": 3,
			"pay_method": "balance",
			"principal_amount": 0,
			"reward_amount": 20,
			"consume_amount": 80,
			"consume_store_id": 1,
			"beneficial_store_id": 5,
			"recharge_order_id": 1670710200224452610,
			"membership_id": 1
		}
	]`
	require.JSONEq(t, expectedJson, string(jsonRes))

	// 校验 membership_balances的值是否正确
	query = `select current_balance_amount, current_principal_amount, current_reward_amount, total_recharge_amount, total_recharge_principal_amount, total_reward_amount
		from membership_balances where membership_id = $1`

	type RechargeInfoResult struct {
		CurrentBalanceAmount         int64 `json:"current_balance_amount"`
		CurrentPrincipalAmount       int64 `json:"current_principal_amount"`
		CurrentRewardAmount          int64 `json:"current_reward_amount"`
		TotalRechargeAmount          int64 `json:"total_recharge_amount"`
		TotalRechargePrincipalAmount int64 `json:"total_recharge_princiapl_amount"`
		TotalRewardAmount            int64 `json:"total_reward_amount"`
	}
	var currentBalance RechargeInfoResult
	err = pgxscan.Get(ctx, v.pool, &currentBalance, query, membershipID)
	require.Nil(t, err)
	jsonRes, err = json.Marshal(currentBalance)
	require.Nil(t, err)
	expectedJson = `
	{
		"current_balance_amount": 120,
		"current_principal_amount": 60,
		"current_reward_amount": 60,
		"total_recharge_amount": 200,
		"total_recharge_princiapl_amount": 100,
		"total_reward_amount": 100
	}`
	require.JSONEq(t, expectedJson, string(jsonRes))

}
