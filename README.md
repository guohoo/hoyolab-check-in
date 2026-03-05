# HoYoLAB Daliy Check-In Helper
目前仅支持国际服签到，包含原神、崩坏3、崩坏：星穹铁道、未定事件簿、绝区零😋。

## 使用方法
Github Actions 中设置 `USER_COOKIE` 和 `ENABLED_GAMES` 两个环境变量（secret），格式如下。

```
// 示例
USER_COOKIE      {"用户名1": "cookie1", "用户名2":"cookie2"}
ENABLED_GAMES    {"用户名1": ["gi", "hk3", "hkrpg", "nxx", "zzz"], "用户名2": ["gi", "hk3", "zzz"]}
```

补充说明：
1. USER_COOKIE 存储用户 cookie，可以在签到网页抓取，一定要包含 ltoken_v2 和 ltuid_v2 这两个字段。
2. ENABLED_GAMES 存储需要启用签到的游戏，其中 "gi"、"hk3"、"hkrpg"、"nxx"、"zzz" 依次对应原神、崩坏3、崩坏：星穹铁道、未定事件簿、绝区零。
3. 多用户使用 "," 分隔。

## 输出效果
```
🚀 开始为 1 个账户执行自动签到任务...
✔️ [Mashiro] Genshin: OK
✔️ [Mashiro] Honkai_3: OK
❌ [Mashiro] Star_Rail: 开拓者，你已经签到过了~
❌ [Mashiro] Tears_of_Themis: 未检测到在游戏内创建角色
❌ [Mashiro] Zenless_Zone_Zero: 绳匠，你已经签到过了~
✨ 签到任务完成 ~
```

## TODO
- 签到结果提醒（Telegram、Discord等）
