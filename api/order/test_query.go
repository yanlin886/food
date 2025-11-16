package order

const (
	queryCase1_1 = `update memberships set entry_store_id = null where id = 1`

	// 重置产品+优惠券信息
	rebuildProductCouponQuery = `
		truncate table skus cascade;
		truncate products cascade;
		truncate category_products cascade;
		truncate coupons cascade;
		truncate coupon_sends cascade;
		truncate memberships cascade;
		truncate membership_membership_card_level_relations cascade;
		truncate membership_coupons cascade;
		truncate orders cascade;

		INSERT INTO category_products(id, tenant_id, name, image, sort, is_sale, created_at, updated_at) values
    		(1, 'yanlin', '面1', 'xx.png', 1, true, '2023-12-12 12:30:30', '2023-12-12 12:30:30'),
    		(2, 'yanlin', '面2', 'xx.png', 1, true, '2023-12-12 12:30:30', '2023-12-12 12:30:30'),
    		(3, 'yanlin', '面3', 'xx.png', 1, true, '2023-12-12 12:30:30', '2023-12-12 12:30:30'),
    		(4, 'yanlin', '面4', 'xx.png', 1, true, '2023-12-12 12:30:30', '2023-12-12 12:30:30'),
    		(5, 'yanlin', '面5', 'xx.png', 1, true, '2023-12-12 12:30:30', '2023-12-12 12:30:30'),
    		(6, 'yanlin', '面6', 'xx.png', 1, true, '2023-12-12 12:30:30', '2023-12-12 12:30:30');

		INSERT INTO products(id, tenant_id, category_id, title, image, is_sale, sku_data, unit, amount, created_at, updated_at) values
    		(1, 'yanlin', 1, '红烧牛肉面1', 'xx.png', false, '{"attr_list": [{"name": "大份", "parent": "份"}, {"name": "中份", "parent": "份"}, {"name": "小份", "parent": "份"}], "spec_list": [{"name": "份"}]}', '碗', '1200', '2022-12-12 12:20:30', '2022-12-12 12:30:30'),
    		(2, 'yanlin', 2, '红烧牛肉面2', 'xx.png', true, '{"attr_list": [{"name": "大份", "parent": "份"}, {"name": "中份", "parent": "份"}, {"name": "小份", "parent": "份"}], "spec_list": [{"name": "份"}]}', '碗', '1200', '2022-12-12 12:20:30', '2022-12-12 12:30:30'),
    		(3, 'yanlin', 3, '红烧牛肉面3', 'xx.png', true, '{"attr_list": [{"name": "大份", "parent": "份"}, {"name": "中份", "parent": "份"}, {"name": "小份", "parent": "份"}], "spec_list": [{"name": "份"}]}', '碗', '1200', '2022-12-12 12:20:30', '2022-12-12 12:30:30'),
    		(4, 'yanlin', 4, '红烧牛肉面4', 'xx.png', true, '{"attr_list": [{"name": "大份", "parent": "份"}, {"name": "中份", "parent": "份"}, {"name": "小份", "parent": "份"}], "spec_list": [{"name": "份"}]}', '碗', '1200', '2022-12-12 12:20:30', '2022-12-12 12:30:30'),
    		(5, 'yanlin', 5, '红烧牛肉面5', 'xx.png', true, '{"attr_list": [{"name": "大份", "parent": "份"}, {"name": "中份", "parent": "份"}, {"name": "小份", "parent": "份"}], "spec_list": [{"name": "份"}]}', '碗', '1200', '2022-12-12 12:20:30', '2022-12-12 12:30:30'),
    		(6, 'yanlin', 1, '红烧牛肉面6', 'xx.png', true, '{"attr_list": [{"name": "大份", "parent": "份"}, {"name": "中份", "parent": "份"}, {"name": "小份", "parent": "份"}], "spec_list": [{"name": "份"}]}', '碗', '1200', '2022-12-12 12:20:30', '2022-12-12 12:30:30'),
    		(7, 'yanlin', 2, '红烧牛肉面7', 'xx.png', true, '{"attr_list": [{"name": "大份", "parent": "份"}, {"name": "中份", "parent": "份"}, {"name": "小份", "parent": "份"}], "spec_list": [{"name": "份"}]}', '碗', '1200', '2022-12-12 12:20:30', '2022-12-12 12:30:30');

		INSERT INTO skus(id, tenant_id, product_id, name, amount) values
    		(1, 'yanlin', 1, '大份', 2000),
    		(2, 'yanlin', 1, '中份', 1800),
    		(3, 'yanlin', 1, '小份', 1600),
    		(4, 'yanlin', 2, '大份', 2000),
    		(5, 'yanlin', 2, '中份', 1800),
    		(6, 'yanlin', 2, '小份', 1600),
    		(7, 'yanlin', 3, '大份', 2000),
    		(8, 'yanlin', 3, '中份', 1800),
    		(9, 'yanlin', 3, '小份', 1600),
    		(10, 'yanlin', 4, '大份', 2000),
    		(11, 'yanlin', 4, '中份', 1800),
    		(12, 'yanlin', 4, '小份', 1600),
    		(13, 'yanlin', 5, '大份', 2000),
    		(14, 'yanlin', 5, '中份', 1800),
    		(15, 'yanlin', 5, '小份', 1600),
    		(16, 'yanlin', 6, '大份', 2000),
    		(17, 'yanlin', 6, '中份', 1800),
    		(18, 'yanlin', 6, '小份', 1600),
    		(19, 'yanlin', 7, '大份', 2000),
    		(20, 'yanlin', 7, '中份', 1800),
    		(21, 'yanlin', 7, '小份', 1600);

		INSERT INTO coupons (id, tenant_id, name, category, base_config, usage_config, sent_number, used_number, created_at, updated_at) VALUES
			(1, 'yanlin', '优惠券222', 'money_reduce', '{"inner_note": "test", "validation": {"category": "fixed", "fixed_end_at": "2023-06-29 18:30:03", "fixed_started_at": "2023-06-04 18:30:01", "relative_delay_days": 0, "relative_valid_days": 22}, "reduce_amount": "0"}', '{"door_money": "0", "usage_note": "饿死的", "is_category": false, "not_use_day": {"not_use_day": [2], "legal_holidays_not_use": true}, "product_ids": ["1", "2"], "can_use_time": [{"end_time": "09:00:00", "start_time": "03:00:00"}, {"end_time": "07:00:00", "start_time": "01:00:00"}], "category_ids": ["1", "2"], "exchange_sort": "asc", "is_applicable": false, "is_all_product": false, "is_membership_overlay": true}', 0, 0, '2023-06-19 10:30:24.364515', '2023-06-19 10:30:24.364515'),
			(2, 'yanlin', '优惠券111', 'product_exchange', '{"inner_note": "三方", "validation": {"category": "fixed", "fixed_end_at": "2023-06-23 18:26:43", "fixed_started_at": "2023-06-12 18:26:42", "relative_delay_days": 0, "relative_valid_days": 0}, "reduce_amount": "0"}', '{"door_money": "0", "usage_note": "呃呃", "is_category": false, "not_use_day": {"not_use_day": [2, 3, 0], "legal_holidays_not_use": true}, "product_ids": ["1", "2"], "can_use_time": [{"end_time": "07:00:00", "start_time": "03:00:00"}, {"end_time": "09:00:00", "start_time": "03:00:00"}], "category_ids": ["1", "2"], "exchange_sort": "asc", "is_applicable": false, "is_all_product": false, "is_membership_overlay": true}', 0, 0, '2023-06-19 10:27:04.781284', '2023-06-19 10:27:04.781284'),
			(3, 'yanlin', '测试优惠券333', 'money_reduce', '{"inner_note": "内部备注222", "validation": {"category": "", "fixed_end_at": "2023-06-20 14:28:42", "fixed_started_at": "2023-06-20 14:28:42", "relative_delay_days": 0, "relative_valid_days": 0}, "reduce_amount": "2"}', '{"door_money": "3", "usage_note": "使用说明111", "is_category": false, "not_use_day": {"not_use_day": [3, 4, 5, 0], "legal_holidays_not_use": true}, "product_ids": ["3"], "can_use_time": [{"end_time": "06:00:00", "start_time": "03:03:00"}, {"end_time": "12:08:00", "start_time": "06:00:00"}], "category_ids": ["3"], "exchange_sort": "asc", "is_applicable": false, "is_all_product": false, "is_membership_overlay": false}', 0, 0, '2023-06-20 06:01:10.383626', '2023-06-20 06:01:10.383626');

		INSERT INTO coupon_sends (id, tenant_id, name, config, next_send_time, category, status, created_at, updated_at,author_id,author_name) values
			(1, 'yanlin', '情人节活动', '{}', '2022-12-12 12:30:30', 'send', 'active', '2022-12-12 12:30:30', '2099-12-12 12:30:30',1,'jack');

		INSERT INTO memberships ("id", "tenant_id", "openid", "name", "is_enabled", "sex", "phone", "birthday", "age", "created_at", "updated_at", "share_membership_store_group_id") VALUES
			(1, 'yanlin', '1', 'jack', 't', 'male', '1994444', '2023-06-11', 12, '2023-06-11 05:01:02.166868', '2023-06-11 05:01:02.166868', 1),
			(2, 'yanlin', '', 'tom', 'f', 'male', '15051552229', '2023-06-15', 44, '2023-06-13 03:04:05.776414', '2023-06-13 03:04:05.776414', 1);

		INSERT INTO membership_points (membership_id, points, tenant_id, total_points) VALUES
			(1, 0, 'yanlin', 0),
			(2, 0, 'yanlin', 0);

		INSERT INTO membership_coupons (id, membership_id, tenant_id, coupon_id, number, start_at, expire_at, coupon_send_id, coupon_name, coupon_category, coupon_base_config, coupon_usage_config, created_at, updated_at) values
			(1, 1, 'yanlin', 1, 2, '2022-12-12 12:30:30', '2099-12-12 12:30:30', 1, '满200减50券', 'money_reduce', '{}', '{}', '2022-12-12 12:30:30', '2099-12-12 12:30:30'),
			(2, 1, 'yanlin', 2, 1, '2022-12-13 12:30:30', '2098-12-12 12:30:30', 1, '优惠券2', 'money_reduce', '{}', '{}', '2022-12-12 12:30:30', '2099-12-12 12:30:30'),
			(3, 1, 'yanlin', 3, 4, '2022-12-13 12:30:30', '2098-12-12 12:30:30', 1, '优惠券2', 'product_exchange', '{}', '{}', '2022-12-12 12:30:30', '2099-12-12 12:30:30'),
			(4, 2, 'yanlin', 1, 5, '2022-12-13 12:30:30', '2098-12-12 12:30:30', 1, '优惠券2', 'money_reduce', '{}', '{}', '2022-12-12 12:30:30', '2099-12-12 12:30:30');
		
		INSERT INTO membership_membership_card_level_relations ("id", "tenant_id", "membership_card_level_id", "membership_id", "created_at") VALUES
			(1, 'yanlin', 1, 1, '2023-06-14 07:45:49.643824'),
			(2, 'yanlin', 1, 2, '2023-06-14 07:45:49.643824');

		INSERT INTO membership_balances (tenant_id, membership_id, current_balance_amount, current_principal_amount, current_reward_amount, total_recharge_amount, total_recharge_principal_amount, total_reward_amount, created_at, updated_at) VALUES 
			('yanlin', 1, 0, 0, 0, 0, 0, 0, '2023-06-13 03:04:05.776414', '2023-06-13 03:04:05.776414'),		                                                                                                                                                                                                                                                          
			('yanlin', 2, 0, 0, 0, 0, 0, 0, '2023-06-13 03:04:05.776414', '2023-06-13 03:04:05.776414');		                                                                                                                                                                                                                                                          
`

	addUpgradeCaseInitSql = `
		truncate table membership_card_levels cascade;
		INSERT INTO membership_card_levels ("id", "tenant_id", "membership_card_id", "name", "upgrade_condition_config", "member_price_config", "points_config", "level") VALUES 
			(1, 'yanlin', 3, 'VIP1', '{"points_count": 0, "enough_amount": "", "points_enable": false, "stored_amount": "", "consumption_count": 0, "consumption_amount": "", "stored_value_enable": false, "consumption_count_enable": false, "consumption_amount_enable": false}', '{}', '{"obtain": true, "cash_out": true, "multiple": true, "obtain_config": {"obtain_point": 2, "consumption_amount": "1"}, "cash_out_config": {"integral_count": 5, "cash_out_amount": "1", "once_upper_limit": "7", "is_enable_upper_limit": true}, "multiple_config": {"vip_day_enable": true, "vip_day_multiple": 4, "birthday_time_slot": "WEEK", "birthday_time_slot_enable": true, "birthday_time_slot_multiple": 3}}', 10),
			(2,'yanlin',5,'VIP1','{"points_count": 0, "enough_amount": "", "points_enable": false, "stored_amount": "", "consumption_count": 0, "consumption_amount": "", "stored_value_enable": false, "consumption_count_enable": false, "consumption_amount_enable": false}','{}','{"obtain": true, "cash_out": true, "multiple": true, "obtain_config": {"obtain_point": 2, "consumption_amount": "1"}, "cash_out_config": {"integral_count": 0, "cash_out_amount": "", "once_upper_limit": "", "is_enable_upper_limit": false}, "multiple_config": {"vip_day_enable": true, "vip_day_multiple": 4, "birthday_time_slot": "DAY", "birthday_time_slot_enable": true, "birthday_time_slot_multiple": 3}}',0),
			(3, 'yanlin', 6, 'VIP1', '{"points_count": 0, "enough_amount": "", "points_enable": false, "stored_amount": "", "consumption_count": 0, "consumption_amount": "", "stored_value_enable": false, "consumption_count_enable": false, "consumption_amount_enable": false}', '{}', '{}', 0),
			(4, 'yanlin' , 5, '等级1', '{"points_enable":true,"points_count":1,"stored_value_enable":false,"stored_amount":"0.01","consumption_count_enable":false,"consumption_count":1,"consumption_amount_enable":false,"consumption_amount":"0.01","enough_amount":"0"}', '{}', '{"obtain": true, "cash_out": true, "multiple": true, "obtain_config": {"obtain_point": 2, "consumption_amount": "1"}, "cash_out_config": {"integral_count": 0, "cash_out_amount": "0", "once_upper_limit": "0", "is_enable_upper_limit": false}, "multiple_config": {"vip_day_enable": true, "vip_day_multiple": 4, "birthday_time_slot": "DAY", "birthday_time_slot_enable": true, "birthday_time_slot_multiple": 3}}', 1),
			(5, 'yanlin', 5, '等级2', '{"points_enable":false,"points_count":2,"stored_value_enable":true,"stored_amount":"0.01","consumption_count_enable":false,"consumption_count":1,"consumption_amount_enable":false,"consumption_amount":"0.01","enough_amount":"0"}', '{}', '{"obtain": true, "cash_out": true, "multiple": true, "obtain_config": {"obtain_point": 2, "consumption_amount": "1"}, "cash_out_config": {"integral_count": 0, "cash_out_amount": "0", "once_upper_limit": "0", "is_enable_upper_limit": false}, "multiple_config": {"vip_day_enable": true, "vip_day_multiple": 4, "birthday_time_slot": "DAY", "birthday_time_slot_enable": true, "birthday_time_slot_multiple": 3}}', 2),
			(6, 'yanlin', 5, '等级3', '{"points_enable":false,"points_count":2,"stored_value_enable":false,"stored_amount":"0.01","consumption_count_enable":true,"consumption_count":1,"consumption_amount_enable":false,"consumption_amount":"0.01","enough_amount":"0"}', '{}', '{"obtain": true, "cash_out": true, "multiple": true, "obtain_config": {"obtain_point": 2, "consumption_amount": "1"}, "cash_out_config": {"integral_count": 0, "cash_out_amount": "0", "once_upper_limit": "0", "is_enable_upper_limit": false}, "multiple_config": {"vip_day_enable": true, "vip_day_multiple": 4, "birthday_time_slot": "DAY", "birthday_time_slot_enable": true, "birthday_time_slot_multiple": 3}}', 3),
			(7, 'yanlin', 5, '等级4', '{"points_enable":false,"points_count":2,"stored_value_enable":false,"stored_amount":"0.01","consumption_count_enable":false,"consumption_count":1,"consumption_amount_enable":true,"consumption_amount":"0.01","enough_amount":"0"}', '{}', '{"obtain": true, "cash_out": true, "multiple": true, "obtain_config": {"obtain_point": 2, "consumption_amount": "1"}, "cash_out_config": {"integral_count": 0, "cash_out_amount": "0", "once_upper_limit": "0", "is_enable_upper_limit": false}, "multiple_config": {"vip_day_enable": true, "vip_day_multiple": 4, "birthday_time_slot": "DAY", "birthday_time_slot_enable": true, "birthday_time_slot_multiple": 3}}', 4);
				
		truncate table membership_membership_card_level_relations cascade;
		INSERT INTO membership_membership_card_level_relations ("id", "tenant_id", "membership_card_level_id", "membership_id", "created_at") VALUES
			(1, 'yanlin', 1, 5, '2023-06-13 03:04:05.776414'),
			(2, 'yanlin', 1, 12, '2023-06-14 07:45:49.643824'),
			(3, 'yanlin', 1, 1, '2023-06-14 07:45:49.643824');
`
	WechatOrderList = `
		truncate table orders cascade;
		truncate table order_items cascade;
		truncate table order_stores cascade;
		truncate table payments cascade ;

		INSERT INTO orders (id, tenant_id, total_amount, actual_amount, type, status, category, source, created_at, updated_at) VALUES
		       (1668915764205195261, 'yanlin', 3100, 3100, 'recharge', 'canceled', 'in_store_dining', 'wechat', '2023-07-12 17:38:19.801796', '2023-06-14 17:38:19.801796'),
		       (1668915764205195262, 'yanlin', 3300, 1100, 'product', 'completed', 'in_store_dining', 'wechat', '2023-09-11 17:38:19.801796', '2023-06-14 17:38:19.801796'),
		       (1668915764205195263, 'yanlin', 1200, 1200, 'product', 'completed', 'in_store_dining', 'wechat', '2022-06-14 17:38:19.801796', '2023-06-14 17:38:19.801796'),
		       (1668915764205195264, 'yanlin', 6600, 3700, 'product', 'completed', 'in_store_dining', 'web', '2023-01-22 17:38:19.801796', '2023-06-14 17:38:19.801796');

		INSERT INTO order_items (id, tenant_id, order_id, total_amount, unit_amount, quantity, title, image, attributes, created_at) VALUES 
				(1, 'yanlin', 1668915764205195261, 1100, 1100, 1000, '充值31元 送41元 送13积分', NULL, 'https://www.jeck.tang.com/003.gif', '2023-06-14 17:38:19.801796'),
				(2, 'yanlin', 1668915764205195262, 1100, 1100, 1000, '红烧牛肉面', 'https://www.jeck.tang.com/003.gif', NULL, '2023-06-14 17:38:19.801796'),
				(3, 'yanlin', 1668915764205195262, 1200, 1200, 1000, '蛋炒饭', 'https://www.jeck.tang.com/003.gif', NULL, '2023-06-14 17:38:19.801796'),
				(4, 'yanlin', 1668915764205195262, 1000, 1000, 1000, '红烧牛肉面', 'https://www.jeck.tang.com/003.gif', NULL, '2023-06-14 17:38:19.801796'),
				(5, 'yanlin', 1668915764205195263, 1200, 1200, 1000, '红烧牛肉面', 'https://www.jeck.tang.com/003.gif', NULL, '2023-06-14 17:38:19.801796'),
				(6, 'yanlin', 1668915764205195264, 6600, 3700, 1000, '蛋炒饭', 'https://www.jeck.tang.com/003.gif', NULL, '2023-06-14 17:38:19.801796');


		INSERT INTO order_item_products (order_item_id, product_id, sku_id, tenant_id, created_at) VALUES
			(2, 1, 1, 'yanlin', '2023-06-14 17:38:19.801796'),
			(3, 2, 4, 'yanlin', '2023-06-14 17:38:19.801796'),
			(4, 1, 2, 'yanlin', '2023-06-14 17:38:19.801796'),
			(5, 1, 3, 'yanlin', '2023-06-14 17:38:19.801796'),
			(6, 2, 5, 'yanlin', '2023-06-14 17:38:19.801796');

		INSERT INTO order_memberships ("tenant_id", "order_id", "membership_id", "created_at") VALUES 
				('yanlin', 1668915764205195261, 1, '2023-06-16 08:50:41.577945'),
				('yanlin', 1668915764205195262, 1, '2023-06-16 08:50:41.577945'),
				('yanlin', 1668915764205195263, 1, '2023-06-16 08:50:41.577945'),
				('yanlin', 1668915764205195264, 1, '2023-06-16 08:50:41.577945');

		INSERT INTO order_stores ("tenant_id", "order_id", "store_id", "created_at") VALUES 
				('yanlin', 1668915764205195261, 1, '2023-07-28 11:56:21'),
				('yanlin', 1668915764205195262, 1, '2023-07-28 11:56:21'),
				('yanlin', 1668915764205195263, 1, '2023-07-28 11:56:21'),
				('yanlin', 1668915764205195264, 1, '2023-07-28 11:56:21');

		INSERT INTO payments (id, tenant_id, method, order_id, wechat_transaction, amount, status, created_at, updated_at) VALUES 
				(1, 'yanlin', 'wechat', 1668915764205195261, '{}', 3100, 'pending', '2023-07-28 11:56:21', '2023-07-28 11:56:21'),
				(2, 'yanlin', 'wechat', 1668915764205195262, '{}', 3300, 'paid', '2023-07-28 11:56:21', '2023-07-28 11:56:21'),
				(3, 'yanlin', 'balance', 1668915764205195263, '{}', 1200, 'paid', '2023-07-28 11:56:21', '2023-07-28 11:56:21'),
				(4, 'yanlin', 'cash', 1668915764205195264, '{}', 6600, 'paid', '2023-07-28 11:56:21', '2023-07-28 11:56:21');
`

	orderNotSerialNumberSql = `
		truncate table order_serial_numbers cascade;
		truncate order_tables cascade;
		truncate order_service_fees cascade;

		insert into order_tables(id,tenant_id,order_id,table_id,name) values 
		  	(1,'yanlin',1670710200224452609,null,'测试桌台');

		insert into order_service_fees(id,tenant_id,order_id,service_fee_id,name,amount) values 
			(1,'yanlin',1670710200224452609,null,'服务费1',10000)
	`
	wechatBalanceListSql = `
		truncate table orders cascade;
		truncate table order_items cascade;
		truncate table order_memberships cascade;

		INSERT INTO orders ("id", "tenant_id", "total_amount", "actual_amount", "type", "status", "created_at", "updated_at") VALUES
			(1668915764205195261, 'yanlin', 3100, 3100, 'recharge', 'canceled', '2023-07-12 17:38:19.801796', '2023-06-14 17:38:19.801796'),
			(1668915764205195262, 'yanlin', 3300, 1100, 'product', 'completed', '2023-09-11 17:38:19.801796', '2023-06-14 17:38:19.801796'),
			(1668915764205195263, 'yanlin', 1200, 1200, 'product', 'completed', '2022-06-14 17:38:19.801796', '2023-06-14 17:38:19.801796'),
			(1668915764205195264, 'yanlin', 6600, 3700, 'product', 'completed', '2023-01-22 17:38:19.801796', '2023-06-14 17:38:19.801796'),
			(1668915764205195265, 'yanlin', 4100, 3100, 'recharge', 'completed', '2023-07-12 17:38:19.801796', '2023-06-14 17:38:19.801796');

		INSERT INTO order_items ("id", "tenant_id", "order_id", "total_amount", "unit_amount", "quantity", "title", "image", "attributes", "created_at") VALUES 
			(1, 'yanlin', 1668915764205195261, 1100, 1100, 1, '充值31元 送41元 送13积分', NULL, 'https://www.jeck.tang.com/003.gif', '2023-06-14 17:38:19.801796'),
			(2, 'yanlin', 1668915764205195262, 1100, 1100, 1, '青菜面', NULL, 'https://www.jeck.tang.com/003.gif', '2023-06-14 17:38:19.801796'),
			(3, 'yanlin', 1668915764205195262, 1200, 1200, 1, '牛肉面', NULL, 'https://www.jeck.tang.com/003.gif', '2023-06-14 17:38:19.801796'),
			(4, 'yanlin', 1668915764205195262, 1000, 1000, 1, '白水面', NULL, 'https://www.jeck.tang.com/003.gif', '2023-06-14 17:38:19.801796'),
			(5, 'yanlin', 1668915764205195263, 1200, 1200, 1, '蛋炒饭', NULL, 'https://www.jeck.tang.com/003.gif', '2023-06-14 17:38:19.801796'),
			(6, 'yanlin', 1668915764205195264, 6600, 3700, 1, '清锅', NULL, 'https://www.jeck.tang.com/003.gif', '2023-06-14 17:38:19.801796'),
			(7, 'yanlin', 1668915764205195265, 4100, 3100, 1, '充值31元 送41元 送13积分', NULL, 'https://www.jeck.tang.com/003.gif', '2023-06-14 17:38:19.801796');

		INSERT INTO order_memberships ("tenant_id", "order_id", "membership_id", "created_at") VALUES 
			('yanlin', 1668915764205195261, 1, '2023-06-16 08:50:41.577945'),
			('yanlin', 1668915764205195262, 5, '2023-06-16 08:50:41.577945'),
			('yanlin', 1668915764205195263, 1, '2023-06-16 08:50:41.577945'),
			('yanlin', 1668915764205195264, 1, '2023-06-16 08:50:41.577945'),
			('yanlin', 1668915764205195265, 1, '2023-06-16 08:50:41.577945');

		INSERT INTO payments ("id", "tenant_id", "method", "order_id", "wechat_transaction", "amount", "status", "created_at", "updated_at") VALUES
			(1, 'yanlin', 'cash', 1668915764205195263, '{}', 123, 'paid', '2023-07-31 16:41:16', '2023-07-31 16:41:18'),
			(2, 'yanlin', 'balance', 1668915764205195264, '{}', 123, 'paid', '2023-07-31 16:41:16', '2023-07-31 16:41:18'),
			(3, 'yanlin', 'cash', 1668915764205195265, '{}', 123, 'paid', '2023-07-31 16:41:16', '2023-07-31 16:41:18');
			
`

	orderCopyTestSQL = `
		truncate table products cascade;
		truncate table skus cascade;
		truncate table order_memberships cascade;
		truncate table order_item_products cascade;
		truncate table order_item_ingredients cascade;
		truncate table order_items cascade;
		truncate table orders cascade;
		truncate table menu_products cascade;
		truncate table menu_product_skus cascade;
		truncate table menu_product_ingredient_options cascade;
		truncate table menu_product_ingredients cascade;

		INSERT INTO products(id, tenant_id, category_id, title, image, is_sale, sku_data, unit, amount, created_at, updated_at) values
			(1, 'yanlin', 1, '红烧牛肉面', 'xx.png', true, '{"attr_list": [{"name": "大份", "parent": "份"}, {"name": "中份", "parent": "份"}], "spec_list": [{"name": "份"}]}', '碗', '1200', '2022-12-12 12:20:30', '2022-12-12 12:30:30'),
			(2, 'yanlin', 2, '蛋炒饭', NULL, 't', '{}', '', 11, '2023-06-09 15:18:09.726443', '2023-06-09 15:18:11.638006');

		INSERT INTO skus(id, tenant_id, product_id, name, amount) values
			(1, 'yanlin', 1, '大份', 2000),
			(2, 'yanlin', 1, '中份', 1800),
			(3, 'yanlin', 2, '小份', 1600);
		
		INSERT INTO menu_products ( "menu_id", "product_id", "tenant_id", "created_at", "updated_at") VALUES
			( 1, 1, 'yanlin', '2023-07-22 14:04:47', '2023-07-22 14:04:50');

		INSERT INTO menu_product_skus ("id", "tenant_id", "menu_id", "product_id", "sku_id", "regular_amount", "member_amount","pack_amount", "created_at") VALUES
			(1, 'yanlin', 1, 1, 2, 124, 123,100, '2023-07-12 09:56:25');

		INSERT INTO menu_product_ingredients ("menu_id", "product_id","ingredient_id", "choose_config", "created_at", "updated_at", "tenant_id") VALUES
			(1, 1, 1, '{"is_required": true, "is_most_required": true, "required_choose_count": 2, "most_choose_count": 3}', '2023-07-14 11:34:35', '2023-07-14 11:34:38', 'yanlin');
		
		INSERT INTO menu_product_ingredient_options ("tenant_id","menu_id","product_id", "ingredient_id", "ingredient_option_id", "amount", "created_at") VALUES
			('yanlin', 1, 1, 1, 1, 321, '2023-07-20 11:05:06'),
			('yanlin', 1, 1, 1, 2, 322, '2023-07-20 11:05:06');

		INSERT INTO orders ("id", "tenant_id", "total_amount", "actual_amount", "type", "status","note", "created_at", "updated_at") VALUES
    		(1670710200224452608, 'yanlin', 3, 3, 'product', 'pending_delivery', 'note1','2023-06-19 16:28:46.700739', '2023-06-19 16:28:46.700739'),
    		(1670710200224452609, 'yanlin', 3, 3, 'product', 'pending_delivery', 'note1','2023-06-19 16:28:46.700739', '2023-06-19 16:28:46.700739');

		INSERT INTO order_stores ("tenant_id", "order_id", "store_id", "created_at") VALUES 
			('yanlin', 1670710200224452608, 1, '2023-07-28 11:56:21'),
			('yanlin', 1670710200224452609, 1, '2023-07-28 11:56:21');

		INSERT INTO order_memberships ("tenant_id", "order_id", "membership_id", "created_at") VALUES
			('yanlin', 1670710200224452608, 5, '2023-06-19 08:28:46.684328'),
			('yanlin', 1670710200224452609, 5, '2023-06-19 08:28:46.684328');

		INSERT INTO order_items ("id", "tenant_id", "order_id", "total_amount", "original_total_amount", "unit_amount", "quantity", "title", "image", "attributes", "created_at") VALUES
			(1, 'yanlin', 1670710200224452608, 3, 5, 3, 1000, '红烧牛肉面', NULL, '中份', '2023-06-19 16:28:46.700739'),
			(2, 'yanlin', 1670710200224452608, 3, 5, 3, 1000, '荤菜', NULL, '鸡蛋', '2023-06-19 16:28:46.700739'),
			(3, 'yanlin', 1670710200224452608, 3, 5, 3, 3000, '荤菜', NULL, '牛肉', '2023-06-19 16:28:46.700739'),
			(4, 'yanlin', 1670710200224452608, 3, 6, 3, 1000, '羊肉面', NULL, '默认', '2023-06-19 16:28:46.700739'),
			(5, 'yanlin', 1670710200224452609, 3, 5, 3, 1000, '红烧牛肉面', NULL, '中份', '2023-06-19 16:28:46.700739');

		insert into order_item_ingredients(order_item_id,ingredient_id,ingredient_option_id,tenant_id,created_at) VALUES
			(2,1,1,'yanlin','2023-06-19 16:28:46.700739'),
			(3,1,2,'yanlin','2023-06-19 16:28:46.700739');

		INSERT INTO order_item_products ("order_item_id", "product_id", "sku_id", "tenant_id", "created_at") VALUES
			(1, 1, 2, 'yanlin', '2023-06-19 08:28:46.684328'),
			(5, 1, 2, 'yanlin', '2023-06-19 08:28:46.684328');

		insert into order_item_self_relations(order_item_id,order_item_parent_id,tenant_id,created_at)VALUES
			(2, 1, 'yanlin','2023-06-19 16:28:46.700739'),
			(3, 1, 'yanlin','2023-06-19 16:28:46.700739');
	`
	orderReverseTestSQL = `
		truncate table products cascade;
		truncate table skus cascade;
		truncate table order_memberships cascade;
		truncate table order_item_products cascade;
		truncate table order_item_ingredients cascade;
		truncate table order_items cascade;
		truncate table orders cascade;
		truncate table menu_products cascade;
		truncate table menu_product_skus cascade;
		truncate table menu_product_ingredient_options cascade;
		truncate table menu_product_ingredients cascade;
		truncate table payments cascade;

		INSERT INTO products(id, tenant_id, category_id, title, image, is_sale, sku_data, unit, amount, created_at, updated_at) values
			(1, 'yanlin', 1, '红烧牛肉面', 'xx.png', true, '{"attr_list": [{"name": "大份", "parent": "份"}, {"name": "中份", "parent": "份"}], "spec_list": [{"name": "份"}]}', '碗', '1200', '2022-12-12 12:20:30', '2022-12-12 12:30:30'),
			(2, 'yanlin', 2, '蛋炒饭', NULL, 't', '{}', '', 11, '2023-06-09 15:18:09.726443', '2023-06-09 15:18:11.638006');

		INSERT INTO skus(id, tenant_id, product_id, name, amount) values
			(1, 'yanlin', 1, '大份', 2000),
			(2, 'yanlin', 1, '中份', 1800),
			(3, 'yanlin', 2, '小份', 1600);
		
		INSERT INTO menu_products ( "menu_id", "product_id", "tenant_id", "created_at", "updated_at") VALUES
			( 1, 1, 'yanlin', '2023-07-22 14:04:47', '2023-07-22 14:04:50');

		INSERT INTO menu_product_skus ("id", "tenant_id", "menu_id", "product_id", "sku_id", "regular_amount", "member_amount","pack_amount", "created_at") VALUES
			(1, 'yanlin', 1, 1, 2, 124, 123,100, '2023-07-12 09:56:25');

		INSERT INTO menu_product_ingredients ("menu_id", "product_id","ingredient_id", "choose_config", "created_at", "updated_at", "tenant_id") VALUES
			(1, 1, 1, '{"is_required": true, "is_most_required": true, "required_choose_count": 2, "most_choose_count": 3}', '2023-07-14 11:34:35', '2023-07-14 11:34:38', 'yanlin');
		
		INSERT INTO menu_product_ingredient_options ("tenant_id","menu_id","product_id", "ingredient_id", "ingredient_option_id", "amount", "created_at") VALUES
			('yanlin', 1, 1, 1, 1, 321, '2023-07-20 11:05:06'),
			('yanlin', 1, 1, 1, 2, 322, '2023-07-20 11:05:06');

		INSERT INTO orders ("id", "tenant_id", "total_amount", "actual_amount", "type", "status","note", "created_at", "updated_at") VALUES
    		(1670710200224452608, 'yanlin', 3, 3, 'product', 'pending_delivery', 'note1','2023-06-19 16:28:46.700739', '2023-06-19 16:28:46.700739'),
    		(1670710200224452609, 'yanlin', 3, 3, 'product', 'pending_delivery', 'note1','2023-06-19 16:28:46.700739', '2023-06-19 16:28:46.700739');

		INSERT INTO order_stores ("tenant_id", "order_id", "store_id", "created_at") VALUES 
			('yanlin', 1670710200224452608, 1, '2023-07-28 11:56:21'),
			('yanlin', 1670710200224452609, 1, '2023-07-28 11:56:21');

		INSERT INTO order_memberships ("tenant_id", "order_id", "membership_id", "created_at") VALUES
			('yanlin', 1670710200224452608, 5, '2023-06-19 08:28:46.684328'),
			('yanlin', 1670710200224452609, 5, '2023-06-19 08:28:46.684328');

		INSERT INTO order_items ("id", "tenant_id", "order_id", "total_amount", "original_total_amount", "unit_amount", "quantity", "title", "image", "attributes", "created_at") VALUES
			(1, 'yanlin', 1670710200224452608, 3, 5, 3, 1000, '红烧牛肉面', NULL, '中份', '2023-06-19 16:28:46.700739'),
			(2, 'yanlin', 1670710200224452608, 3, 5, 3, 1000, '荤菜', NULL, '鸡蛋', '2023-06-19 16:28:46.700739'),
			(3, 'yanlin', 1670710200224452608, 3, 5, 3, 3000, '荤菜', NULL, '牛肉', '2023-06-19 16:28:46.700739'),
			(4, 'yanlin', 1670710200224452608, 3, 6, 3, 1000, '羊肉面', NULL, '默认', '2023-06-19 16:28:46.700739'),
			(5, 'yanlin', 1670710200224452609, 3, 5, 3, 1000, '红烧牛肉面', NULL, '中份', '2023-06-19 16:28:46.700739');

		insert into order_item_ingredients(order_item_id,ingredient_id,ingredient_option_id,tenant_id,created_at) VALUES
			(2,1,1,'yanlin','2023-06-19 16:28:46.700739'),
			(3,1,2,'yanlin','2023-06-19 16:28:46.700739');

		INSERT INTO order_item_products ("order_item_id", "product_id", "sku_id", "tenant_id", "created_at") VALUES
			(1, 1, 2, 'yanlin', '2023-06-19 08:28:46.684328'),
			(5, 1, 2, 'yanlin', '2023-06-19 08:28:46.684328');

		insert into order_item_self_relations(order_item_id,order_item_parent_id,tenant_id,created_at)VALUES
			(2, 1, 'yanlin','2023-06-19 16:28:46.700739'),
			(3, 1, 'yanlin','2023-06-19 16:28:46.700739');

		INSERT INTO payments ("id", "tenant_id", "method", "order_id", "wechat_transaction", "amount", "status", "created_at", "updated_at") VALUES
			(1, 'yanlin', 'balance', 1670710200224452608, NULL, 3, 'paid', '2023-06-19 16:28:46.722471', '2023-06-19 16:28:46.722471');
	`
	orderListTestSQL = `
		truncate table products cascade;
		truncate table skus cascade;
		truncate table order_memberships cascade;
		truncate table order_item_products cascade;
		truncate table order_item_ingredients cascade;
		truncate table order_items cascade;
		truncate table orders cascade;
		truncate table menu_products cascade;
		truncate table menu_product_skus cascade;
		truncate table menu_product_ingredient_options cascade;
		truncate table menu_product_ingredients cascade;
		truncate table payments cascade;

		INSERT INTO products(id, tenant_id, category_id, title, image, is_sale, sku_data, unit, amount, created_at, updated_at) values
			(1, 'yanlin', 1, '红烧牛肉面', 'xx.png', true, '{"attr_list": [{"name": "大份", "parent": "份"}, {"name": "中份", "parent": "份"}], "spec_list": [{"name": "份"}]}', '碗', '1200', '2022-12-12 12:20:30', '2022-12-12 12:30:30'),
			(2, 'yanlin', 2, '蛋炒饭', NULL, 't', '{}', '', 11, '2023-06-09 15:18:09.726443', '2023-06-09 15:18:11.638006');

		INSERT INTO skus(id, tenant_id, product_id, name, amount) values
			(1, 'yanlin', 1, '大份', 2000),
			(2, 'yanlin', 1, '中份', 1800),
			(3, 'yanlin', 2, '小份', 1600);
		
		INSERT INTO menu_products ( "menu_id", "product_id", "tenant_id", "created_at", "updated_at") VALUES
			( 1, 1, 'yanlin', '2023-07-22 14:04:47', '2023-07-22 14:04:50');

		INSERT INTO menu_product_skus ("id", "tenant_id", "menu_id", "product_id", "sku_id", "regular_amount", "member_amount","pack_amount", "created_at") VALUES
			(1, 'yanlin', 1, 1, 2, 124, 123,100, '2023-07-12 09:56:25');

		INSERT INTO menu_product_ingredients ("menu_id", "product_id","ingredient_id", "choose_config", "created_at", "updated_at", "tenant_id") VALUES
			(1, 1, 1, '{"is_required": true, "is_most_required": true, "required_choose_count": 2, "most_choose_count": 3}', '2023-07-14 11:34:35', '2023-07-14 11:34:38', 'yanlin');
		
		INSERT INTO menu_product_ingredient_options ("tenant_id","menu_id","product_id", "ingredient_id", "ingredient_option_id", "amount", "created_at") VALUES
			('yanlin', 1, 1, 1, 1, 321, '2023-07-20 11:05:06'),
			('yanlin', 1, 1, 1, 2, 322, '2023-07-20 11:05:06');

		INSERT INTO orders ("id", "tenant_id", "total_amount", "actual_amount", "type","status","note", "category", "created_at", "updated_at") VALUES
    		(1670710200224452608, 'yanlin', 3, 3, 'product', 'pending_delivery', 'note1','in_store_dining','2023-06-19 16:28:46.700739', '2023-06-19 16:28:46.700739'),
    		(1670710200224452609, 'yanlin', 3, 3, 'product', 'pending_delivery', 'note1','pick_up','2023-06-19 16:28:46.700739', '2023-06-19 16:28:46.700739');

		INSERT INTO order_operations (tenant_id, order_id, author_name, author_type, event, content) VALUES 
		    ('yanlin', 1670710200224452608, 'admin', 'admin', 'partial_refund', ''),
		    ('yanlin', 1670710200224452608, 'admin', 'admin', 'pay_success', ''),
		    ('yanlin', 1670710200224452609, 'admin', 'admin', 'pay_success', ''),
		    ('yanlin', 1670710200224452609, 'admin', 'admin', 'reverse', '');
		                                                                                                                         

		INSERT INTO order_stores ("tenant_id", "order_id", "store_id", "created_at") VALUES 
			('yanlin', 1670710200224452608, 1, '2023-07-28 11:56:21'),
			('yanlin', 1670710200224452609, 1, '2023-07-28 11:56:21');

		INSERT INTO order_memberships ("tenant_id", "order_id", "membership_id", "created_at") VALUES
			('yanlin', 1670710200224452608, 5, '2023-06-19 08:28:46.684328'),
			('yanlin', 1670710200224452609, 5, '2023-06-19 08:28:46.684328');

		INSERT INTO order_items ("id", "tenant_id", "order_id", "total_amount", "original_total_amount", "unit_amount", "quantity", "title", "image", "attributes", "created_at") VALUES
			(1, 'yanlin', 1670710200224452608, 3, 5, 3, 1, '红烧牛肉面', NULL, '中份', '2023-06-19 16:28:46.700739'),
			(2, 'yanlin', 1670710200224452608, 3, 5, 3, 1, '荤菜', NULL, '鸡蛋', '2023-06-19 16:28:46.700739'),
			(3, 'yanlin', 1670710200224452608, 3, 5, 3, 3, '荤菜', NULL, '牛肉', '2023-06-19 16:28:46.700739'),
			(4, 'yanlin', 1670710200224452608, 3, 6, 3, 1, '羊肉面', NULL, '默认', '2023-06-19 16:28:46.700739'),
			(5, 'yanlin', 1670710200224452609, 3, 5, 3, 1, '红烧牛肉面', NULL, '中份', '2023-06-19 16:28:46.700739');

		insert into order_item_ingredients(order_item_id,ingredient_id,ingredient_option_id,tenant_id,created_at) VALUES
			(2,1,1,'yanlin','2023-06-19 16:28:46.700739'),
			(3,1,2,'yanlin','2023-06-19 16:28:46.700739');

		INSERT INTO order_item_products ("order_item_id", "product_id", "sku_id", "tenant_id", "created_at") VALUES
			(1, 1, 2, 'yanlin', '2023-06-19 08:28:46.684328'),
			(5, 1, 2, 'yanlin', '2023-06-19 08:28:46.684328');

		insert into order_item_self_relations(order_item_id,order_item_parent_id,tenant_id,created_at)VALUES
			(2, 1, 'yanlin','2023-06-19 16:28:46.700739'),
			(3, 1, 'yanlin','2023-06-19 16:28:46.700739');

		INSERT INTO payments ("id", "tenant_id", "method", "order_id", "wechat_transaction", "amount", "status", "created_at", "updated_at") VALUES
			(1, 'yanlin', 'balance', 1670710200224452608, NULL, 3, 'paid', '2023-06-19 16:28:46.722471', '2023-06-19 16:28:46.722471');
	`

	webOrderTest = `
		truncate table_regions cascade;
		truncate tables cascade;
		truncate service_fees cascade;
		truncate table_service_fees cascade;
		truncate order_tables cascade;
		truncate order_service_fees cascade;

		insert into table_regions (id, tenant_id, store_id,name,note) values
			(1,'yanlin',1,'区域A','note1'),
			(2,'yanlin',1,'区域B','note2'),
			(3,'yanlin',1,'区域C','note3'),
			(4,'yanlin',2,'区域D','note1');

		insert into tables(id, tenant_id, store_id, region_id, name, status, usage,cart,number, is_bill_printed) values
			(1,'yanlin',1,1,'台1','pending_open','{"number": 3,"membership": {"id": "1","name": "test_membership","phone": "12345678910"},"start_at":"2023-10-18 07:46:46","note": "note_test","open_type": "wechat"}',
				'{"orders":[{"note":"test","status":"took_order", "order_type":"web","category":"product","products":[{"id":"1","type":"order","extra":[{"id":"1","name":"鸡蛋","category":"做法","quantity":"1","unit_price":"2","subtotal_price_cal":2}],"title":"红烧牛肉面","sku_id":"2","status":"waiting","quantity":"1","attribute":"中份","unit_price":"18","uuid":"1715193983451271168","membership_unit_price":"16","use_membership_unit_price":true}],"uuid":"1715193983451271169","total_amount":"20","actual_amount":"18"},{"uuid":"1715208462788464640","order_type":"web","category":"service_fee","service_fees":{"number":3,"compute_configs":[{"mode":"pro_rata","portion":10},{"mode":"pro_rata","portion":5}]},"note":"","total_amount":"","actual_amount":"","status":"took_order"}]}',1,false),
			(2,'yanlin',1,1,'台2','pending_payment','{"number": 3,"membership": {"id": "1","name": "test_membership","phone": "12345678910"},"start_at":"2023-10-18 07:46:46","note": "note_test","open_type": "wechat"}',
				'{"orders":[{"note":"test","status":"pending_take_order", "order_type":"web","category":"product","products":[{"id":"1","type":"order","extra":[{"id":"1","name":"鸡蛋","category":"做法","quantity":"1","unit_price":"2","subtotal_price_cal":2}],"title":"红烧牛肉面","sku_id":"2","status":"serving","quantity":"1","attribute":"中份","unit_price":"18","uuid":"1715193983451271168","membership_unit_price":"16","use_membership_unit_price":true}],"uuid":"1715193983451271169","total_amount":"20","actual_amount":"18"}]}',2,true),
			(3,'yanlin',1,2,'台3','pending_order','{"number": 3,"start_at":"2023-10-18 07:46:46","note": "note_test","open_type": "web"}','{}',3,false),
			(4,'yanlin',2,4,'台1','pending_open','{}','{}',4,false);

		insert into service_fees(id, tenant_id, store_id, name, compute_config, is_enabled) values
			(1,'yanlin',1,'服务费1','{"mode":"pro_rata","portion":10}',true),
			(2,'yanlin',1,'服务费2','{"mode":"pro_rata","portion":5}',true),
			(3,'yanlin',2,'服务费3','{"mode":"fixed_amount","portion":100}',true);

		insert into table_service_fees(id,tenant_id, table_id, service_fee_id) values
			(1,'yanlin',1,1),
			(2,'yanlin',1,2);
	`

	testQuery1_2 = `
		update store_configs set dinner_config = '{"pay_mode": "order_pay","checkout_mode": "pay_first","receive_mode": "manual","max_auto_receive_amount": "500"}' where store_id = 1
		
	`
	testQuery1_3 = `
		update store_configs set dinner_config = '{"pay_mode": "order_pay","checkout_mode": "pay_first","receive_mode": "manual","max_auto_receive_amount": "500"}', is_zero_pay_allowed = false
		
	`

	// 测试删除优惠券回滚预数据
	refundCouponPreQuery = `
		truncate table membership_coupons cascade;
		truncate table coupon_sends cascade;
		truncate table coupons cascade;
		truncate table memberships cascade;	

		update orders set status = 'pending_payment' where id = 1670710200224452608;

		INSERT INTO coupons (id, tenant_id, name, category, base_config, usage_config, sent_number, used_number, created_at, updated_at, deleted_at) VALUES
			(1, 'yanlin', '优惠券1', 'money_reduce', '{}', '{}', 0, 0, '2023-06-19 10:30:24.364515', '2023-06-19 10:30:24.364515', NULL),
			(2, 'yanlin', '优惠券2', 'money_reduce', '{}', '{}', 0, 0, '2023-06-19 10:27:04.781284', '2023-06-19 10:27:04.781284', NULL),
			(3, 'yanlin', '优惠券3', 'product_exchange', '{}', '{}', 0, 0, '2023-06-20 06:01:10.383626', '2023-06-20 06:01:10.383626', NULL),
			(4, 'yanlin', '优惠券4', 'product_exchange', '{}', '{}', 0, 0, '2023-06-20 06:01:10.383626', '2023-06-20 06:01:10.383626', NULL),
			(5, 'yanlin', '优惠券5', 'product_exchange', '{}', '{}', 0, 0, '2023-06-20 06:01:10.383626', '2023-06-20 06:01:10.383626', now() - '2 day'::interval),
			(6, 'yanlin', '优惠券6', 'product_exchange', '{}', '{}', 0, 0, '2023-06-20 06:01:10.383626', '2023-06-20 06:01:10.383626', now() - '2 day'::interval),
			(7, 'yanlin', '优惠券7', 'product_exchange', '{}', '{}', 0, 0, '2023-06-20 06:01:10.383626', '2023-06-20 06:01:10.383626', now() - '2 day'::interval),
			(8, 'yanlin', '优惠券8', 'product_exchange', '{}', '{}', 0, 0, '2023-06-20 06:01:10.383626', '2023-06-20 06:01:10.383626', now() - '2 day'::interval),
			(9, 'yanlin', '优惠券9', 'product_exchange', '{}', '{}', 0, 0, '2023-06-20 06:01:10.383626', '2023-06-20 06:01:10.383626', now() - '2 day'::interval),
			(10, 'yanlin', '优惠券10', 'product_exchange', '{}', '{}', 0, 0, '2023-06-20 06:01:10.383626', '2023-06-20 06:01:10.383626', NULL);

		INSERT INTO memberships ("id", "tenant_id", "openid", "name", "is_enabled", "sex", "phone", "birthday", "age", "created_at", "updated_at","share_membership_store_group_id") VALUES
    		(1, 'yanlin', '1', 'jack', 't', 'male', '1994444', '2023-06-11', 12, '2023-06-11 05:01:02.166868', '2023-06-11 05:01:02.166868', 1);

		INSERT INTO coupon_sends (id, tenant_id, name, config, next_send_time, category, status, created_at, updated_at,author_id,author_name) values
			(1, 'yanlin', '情人节活动', '{}', '2022-12-12 12:30:30', 'send', 'active', '2022-12-12 12:30:30', '2099-12-12 12:30:30',1,'jack');

		INSERT INTO membership_coupons (id, membership_id, tenant_id, coupon_id, number, start_at, expire_at, coupon_send_id, coupon_name, coupon_category, coupon_base_config, coupon_usage_config, created_at, updated_at) values
			(1, 1, 'yanlin', 1, 1, '2022-12-12 12:30:30', '2023-07-30 12:30:30', 1, '优惠券1', 'money_reduce', '{}', '{}', '2022-12-12 12:30:30', '2099-12-12 12:30:30'),
			(2, 1, 'yanlin', 1, 2, '2022-12-12 12:30:30', '2023-07-30 12:30:30', 1, '优惠券1', 'money_reduce', '{}', '{}', '2022-12-12 12:30:30', '2099-12-12 12:30:30');

		INSERT INTO membership_used_coupons (id, membership_id, tenant_id, coupon_id,start_at, expire_at, coupon_send_id, coupon_name, coupon_category, coupon_base_config, coupon_usage_config,origin_created_at, created_at, membership_coupon_id) values
			(1, 1, 'yanlin', 1, '2022-12-12 12:30:30', '2023-07-30 12:30:30', 1, '优惠券1', 'money_reduce', '{}', '{}', '2022-12-12 12:30:30', '2022-12-12 12:30:30', 1),
			(2, 1, 'yanlin', 1, '2022-12-13 12:30:30', '2023-07-30 12:30:30', 1, '优惠券2', 'money_reduce', '{}', '{}', '2022-12-12 12:30:30', '2022-12-12 12:30:30', 2);
		
		INSERT INTO order_coupons ("id", "tenant_id", "membership_id", "coupon_send_id", "coupon_id", "coupon_property", "order_id", "reduce_amount", "sale_amount", "created_at", membership_used_coupon_id) VALUES
			(1, 'yanlin', 1, 1, 1, '{"membership_coupon_id": "1", "reduce_amount": "12.89", "category": "money_reduce", "name": "满50元减10元"}', 1670710200224452608, 2, 9, '2023-06-19 16:28:46.700739', 1),
			(2, 'yanlin', 1, 1, 1, '{"membership_coupon_id": "2", "reduce_amount": "12.89", "category": "money_reduce", "name": "满50元减10元"}', 1670710200224452609, 2, 9, '2023-06-19 16:28:46.700739', 2);
		
		DELETE FROM membership_coupons WHERE id = 1;
	`

	queryTestCase30 = `update orders set status = 'completed' where id = 1670710200224452612`

	queryTestCase33 = `update recharges set deduction_method = 'proportion'`

	queryTestCase35 = `update orders set status = 'pending_delivery' where id = 1670710200224452608`

	queryTestCase11_1 = `update membership_balances set current_balance_amount = 9900, current_principal_amount = 5000, current_reward_amount = 4900 where membership_id = 1`
)
