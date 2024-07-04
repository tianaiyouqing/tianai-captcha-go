package generator

import (
	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
	"github.com/tianaiyouqing/tianai-captcha-go/common"
	"github.com/tianaiyouqing/tianai-captcha-go/common/model"
	"github.com/tianaiyouqing/tianai-captcha-go/resource"
	"image/color"
	"math"
	"math/rand"
	"unicode/utf8"
)

// 一画到十二画的常用汉字
var chinese_chars = []string{"一", "乙", "二", "十", "丁", "厂", "七", "卜", "人", "入", "八", "九", "几", "儿", "了", "力", "乃", "刀", "又", "三", "于", "干", "亏", "士", "工", "土", "才", "寸", "下", "大", "丈", "与", "万", "上", "小", "口", "巾", "山", "千", "乞", "川", "亿", "个", "勺", "久", "凡", "及", "夕", "丸", "么", "广", "亡", "门", "义", "之", "尸", "弓", "己", "已", "子", "卫", "也", "女", "飞", "刃", "习", "叉", "马", "乡", "丰", "王", "井", "开", "夫", "天", "无", "元", "专", "云", "扎", "艺", "木", "五", "支", "厅", "不", "太", "犬", "区", "历", "尤", "友", "匹", "车", "巨", "牙", "屯", "比", "互", "切", "瓦", "止", "少", "日", "中", "冈", "贝", "内", "水", "见", "午", "牛", "手", "毛", "气", "升", "长", "仁", "什", "片", "仆", "化", "仇", "币", "仍", "仅", "斤", "爪", "反", "介", "父", "从", "今", "凶", "分", "乏", "公", "仓", "月", "氏", "勿", "欠", "风", "丹", "匀", "乌", "凤", "勾", "文", "六", "方", "火", "为", "斗", "忆", "订", "计", "户", "认", "心", "尺", "引", "丑", "巴", "孔", "队", "办", "以", "允", "予", "劝", "双", "书", "幻", "廿", "玉", "刊", "示", "末", "未", "击", "打", "巧", "正", "扑", "扒", "功", "扔", "去", "甘", "世", "古", "节", "本", "术", "可", "丙", "左", "厉", "右", "石", "布", "龙", "平", "灭", "轧", "东", "卡", "北", "占", "业", "旧", "帅", "归", "且", "旦", "目", "叶", "甲", "申", "叮", "电", "号", "田", "由", "史", "只", "央", "兄", "叼", "叫", "另", "叨", "叹", "四", "生", "失", "禾", "丘", "付", "仗", "代", "仙", "们", "仪", "白", "仔", "他", "斥", "瓜", "乎", "丛", "令", "用", "甩", "印", "乐", "句", "匆", "册", "犯", "外", "处", "冬", "鸟", "务", "包", "饥", "主", "市", "立", "闪", "兰", "半", "汁", "汇", "头", "汉", "宁", "穴", "它", "讨", "写", "让", "礼", "训", "必", "议", "讯", "记", "永", "司", "尼", "民", "冯", "辽", "奶", "奴", "加", "召", "皮", "边", "发", "孕", "圣", "对", "台", "矛", "纠", "母", "幼", "丝", "式", "刑", "动", "扛", "寺", "吉", "扣", "考", "托", "老", "执", "巩", "圾", "扩", "扫", "地", "扬", "场", "耳", "共", "芒", "亚", "芝", "朽", "朴", "机", "权", "过", "臣", "再", "协", "西", "压", "厌", "在", "有", "百", "存", "而", "页", "匠", "夸", "夺", "灰", "达", "列", "死", "成", "夹", "轨", "邪", "划", "迈", "毕", "至", "此", "贞", "师", "尘", "尖", "劣", "光", "当", "早", "吐", "吓", "虫", "曲", "团", "同", "吊", "吃", "因", "吸", "吗", "屿", "帆", "岁", "回", "岂", "刚", "则", "肉", "网", "年", "朱", "先", "丢", "舌", "竹", "迁", "乔", "伟", "传", "乒", "乓", "休", "伍", "伏", "优", "伐", "延", "件", "任", "伤", "价", "份", "华", "仰", "仿", "伙", "伪", "自", "血", "向", "似", "后", "行", "舟", "全", "会", "杀", "合", "兆", "企", "众", "爷", "伞", "创", "肌", "朵", "杂", "危", "旬", "旨", "负", "各", "名", "多", "争", "色", "壮", "冲", "冰", "庄", "庆", "亦", "刘", "齐", "交", "次", "衣", "产", "决", "充", "妄", "闭", "问", "闯", "羊", "并", "关", "米", "灯", "州", "汗", "污", "江", "池", "汤", "忙", "兴", "宇", "守", "宅", "字", "安", "讲", "军", "许", "论", "农", "讽", "设", "访", "寻", "那", "迅", "尽", "导", "异", "孙", "阵", "阳", "收", "阶", "阴", "防", "奸", "如", "妇", "好", "她", "妈", "戏", "羽", "观", "欢", "买", "红", "纤", "级", "约", "纪", "驰", "巡", "寿", "弄", "麦", "形", "进", "戒", "吞", "远", "违", "运", "扶", "抚", "坛", "技", "坏", "扰", "拒", "找", "批", "扯", "址", "走", "抄", "坝", "贡", "攻", "赤", "折", "抓", "扮", "抢", "孝", "均", "抛", "投", "坟", "抗", "坑", "坊", "抖", "护", "壳", "志", "扭", "块", "声", "把", "报", "却", "劫", "芽", "花", "芹", "芬", "苍", "芳", "严", "芦", "劳", "克", "苏", "杆", "杠", "杜", "材", "村", "杏", "极", "李", "杨", "求", "更", "束", "豆", "两", "丽", "医", "辰", "励", "否", "还", "歼", "来", "连", "步", "坚", "旱", "盯", "呈", "时", "吴", "助", "县", "里", "呆", "园", "旷", "围", "呀", "吨", "足", "邮", "男", "困", "吵", "串", "员", "听", "吩", "吹", "呜", "吧", "吼", "别", "岗", "帐", "财", "针", "钉", "告", "我", "乱", "利", "秃", "秀", "私", "每", "兵", "估", "体", "何", "但", "伸", "作", "伯", "伶", "佣", "低", "你", "住", "位", "伴", "身", "皂", "佛", "近", "彻", "役", "返", "余", "希", "坐", "谷", "妥", "含", "邻", "岔", "肝", "肚", "肠", "龟", "免", "狂", "犹", "角", "删", "条", "卵", "岛", "迎", "饭", "饮", "系", "言", "冻", "状", "亩", "况", "床", "库", "疗", "应", "冷", "这", "序", "辛", "弃", "冶", "忘", "闲", "间", "闷", "判", "灶", "灿", "弟", "汪", "沙", "汽", "沃", "泛", "沟", "没", "沈", "沉", "怀", "忧", "快", "完", "宋", "宏", "牢", "究", "穷", "灾", "良", "证", "启", "评", "补", "初", "社", "识", "诉", "诊", "词", "译", "君", "灵", "即", "层", "尿", "尾", "迟", "局", "改", "张", "忌", "际", "陆", "阿", "陈", "阻", "附", "妙", "妖", "妨", "努", "忍", "劲", "鸡", "驱", "纯", "纱", "纳", "纲", "驳", "纵", "纷", "纸", "纹", "纺", "驴", "纽", "奉", "玩", "环", "武", "青", "责", "现", "表", "规", "抹", "拢", "拔", "拣", "担", "坦", "押", "抽", "拐", "拖", "拍", "者", "顶", "拆", "拥", "抵", "拘", "势", "抱", "垃", "拉", "拦", "拌", "幸", "招", "坡", "披", "拨", "择", "抬", "其", "取", "苦", "若", "茂", "苹", "苗", "英", "范", "直", "茄", "茎", "茅", "林", "枝", "杯", "柜", "析", "板", "松", "枪", "构", "杰", "述", "枕", "丧", "或", "画", "卧", "事", "刺", "枣", "雨", "卖", "矿", "码", "厕", "奔", "奇", "奋", "态", "欧", "垄", "妻", "轰", "顷", "转", "斩", "轮", "软", "到", "非", "叔", "肯", "齿", "些", "虎", "虏", "肾", "贤", "尚", "旺", "具", "果", "味", "昆", "国", "昌", "畅", "明", "易", "昂", "典", "固", "忠", "咐", "呼", "鸣", "咏", "呢", "岸", "岩", "帖", "罗", "帜", "岭", "凯", "败", "贩", "购", "图", "钓", "制", "知", "垂", "牧", "物", "乖", "刮", "秆", "和", "季", "委", "佳", "侍", "供", "使", "例", "版", "侄", "侦", "侧", "凭", "侨", "佩", "货", "依", "的", "迫", "质", "欣", "征", "往", "爬", "彼", "径", "所", "舍", "金", "命", "斧", "爸", "采", "受", "乳", "贪", "念", "贫", "肤", "肺", "肢", "肿", "胀", "朋", "股", "肥", "服", "胁", "周", "昏", "鱼", "兔", "狐", "忽", "狗", "备", "饰", "饱", "饲", "变", "京", "享", "店", "夜", "庙", "府", "底", "剂", "郊", "废", "净", "盲", "放", "刻", "育", "闸", "闹", "郑", "券", "卷", "单", "炒", "炊", "炕", "炎", "炉", "沫", "浅", "法", "泄", "河", "沾", "泪", "油", "泊", "沿", "泡", "注", "泻", "泳", "泥", "沸", "波", "泼", "泽", "治", "怖", "性", "怕", "怜", "怪", "学", "宝", "宗", "定", "宜", "审", "宙", "官", "空", "帘", "实", "试", "郎", "诗", "肩", "房", "诚", "衬", "衫", "视", "话", "诞", "询", "该", "详", "建", "肃", "录", "隶", "居", "届", "刷", "屈", "弦", "承", "孟", "孤", "陕", "降", "限", "妹", "姑", "姐", "姓", "始", "驾", "参", "艰", "线", "练", "组", "细", "驶", "织", "终", "驻", "驼", "绍", "经", "贯", "奏", "春", "帮", "珍", "玻", "毒", "型", "挂", "封", "持", "项", "垮", "挎", "城", "挠", "政", "赴", "赵", "挡", "挺", "括", "拴", "拾", "挑", "指", "垫", "挣", "挤", "拼", "挖", "按", "挥", "挪", "某", "甚", "革", "荐", "巷", "带", "草", "茧", "茶", "荒", "茫", "荡", "荣", "故", "胡", "南", "药", "标", "枯", "柄", "栋", "相", "查", "柏", "柳", "柱", "柿", "栏", "树", "要", "咸", "威", "歪", "研", "砖", "厘", "厚", "砌", "砍", "面", "耐", "耍", "牵", "残", "殃", "轻", "鸦", "皆", "背", "战", "点", "临", "览", "竖", "省", "削", "尝", "是", "盼", "眨", "哄", "显", "哑", "冒", "映", "星", "昨", "畏", "趴", "胃", "贵", "界", "虹", "虾", "蚁", "思", "蚂", "虽", "品", "咽", "骂", "哗", "咱", "响", "哈", "咬", "咳", "哪", "炭", "峡", "罚", "贱", "贴", "骨", "钞", "钟", "钢", "钥", "钩", "卸", "缸", "拜", "看", "矩", "怎", "牲", "选", "适", "秒", "香", "种", "秋", "科", "重", "复", "竿", "段", "便", "俩", "贷", "顺", "修", "保", "促", "侮", "俭", "俗", "俘", "信", "皇", "泉", "鬼", "侵", "追", "俊", "盾", "待", "律", "很", "须", "叙", "剑", "逃", "食", "盆", "胆", "胜", "胞", "胖", "脉", "勉", "狭", "狮", "独", "狡", "狱", "狠", "贸", "怨", "急", "饶", "蚀", "饺", "饼", "弯", "将", "奖", "哀", "亭", "亮", "度", "迹", "庭", "疮", "疯", "疫", "疤", "姿", "亲", "音", "帝", "施", "闻", "阀", "阁", "差", "养", "美", "姜", "叛", "送", "类", "迷", "前", "首", "逆", "总", "炼", "炸", "炮", "烂", "剃", "洁", "洪", "洒", "浇", "浊", "洞", "测", "洗", "活", "派", "洽", "染", "济", "洋", "洲", "浑", "浓", "津", "恒", "恢", "恰", "恼", "恨", "举", "觉", "宣", "室", "宫", "宪", "突", "穿", "窃", "客", "冠", "语", "扁", "袄", "祖", "神", "祝", "误", "诱", "说", "诵", "垦", "退", "既", "屋", "昼", "费", "陡", "眉", "孩", "除", "险", "院", "娃", "姥", "姨", "姻", "娇", "怒", "架", "贺", "盈", "勇", "怠", "柔", "垒", "绑", "绒", "结", "绕", "骄", "绘", "给", "络", "骆", "绝", "绞", "统", "耕", "耗", "艳", "泰", "珠", "班", "素", "蚕", "顽", "盏", "匪", "捞", "栽", "捕", "振", "载", "赶", "起", "盐", "捎", "捏", "埋", "捉", "捆", "捐", "损", "都", "哲", "逝", "捡", "换", "挽", "热", "恐", "壶", "挨", "耻", "耽", "恭", "莲", "莫", "荷", "获", "晋", "恶", "真", "框", "桂", "档", "桐", "株", "桥", "桃", "格", "校", "核", "样", "根", "索", "哥", "速", "逗", "栗", "配", "翅", "辱", "唇", "夏", "础", "破", "原", "套", "逐", "烈", "殊", "顾", "轿", "较", "顿", "毙", "致", "柴", "桌", "虑", "监", "紧", "党", "晒", "眠", "晓", "鸭", "晃", "晌", "晕", "蚊", "哨", "哭", "恩", "唤", "啊", "唉", "罢", "峰", "圆", "贼", "贿", "钱", "钳", "钻", "铁", "铃", "铅", "缺", "氧", "特", "牺", "造", "乘", "敌", "秤", "租", "积", "秧", "秩", "称", "秘", "透", "笔", "笑", "笋", "债", "借", "值", "倚", "倾", "倒", "倘", "俱", "倡", "候", "俯", "倍", "倦", "健", "臭", "射", "躬", "息", "徒", "徐", "舰", "舱", "般", "航", "途", "拿", "爹", "爱", "颂", "翁", "脆", "脂", "胸", "胳", "脏", "胶", "脑", "狸", "狼", "逢", "留", "皱", "饿", "恋", "桨", "浆", "衰", "高", "席", "准", "座", "脊", "症", "病", "疾", "疼", "疲", "效", "离", "唐", "资", "凉", "站", "剖", "竞", "部", "旁", "旅", "畜", "阅", "羞", "瓶", "拳", "粉", "料", "益", "兼", "烤", "烘", "烦", "烧", "烛", "烟", "递", "涛", "浙", "涝", "酒", "涉", "消", "浩", "海", "涂", "浴", "浮", "流", "润", "浪", "浸", "涨", "烫", "涌", "悟", "悄", "悔", "悦", "害", "宽", "家", "宵", "宴", "宾", "窄", "容", "宰", "案", "请", "朗", "诸", "读", "扇", "袜", "袖", "袍", "被", "祥", "课", "谁", "调", "冤", "谅", "谈", "谊", "剥", "恳", "展", "剧", "屑", "弱", "陵", "陶", "陷", "陪", "娱", "娘", "通", "能", "难", "预", "桑", "绢", "绣", "验", "继", "球", "理", "捧", "堵", "描", "域", "掩", "捷", "排", "掉", "堆", "推", "掀", "授", "教", "掏", "掠", "培", "接", "控", "探", "据", "掘", "职", "基", "著", "勒", "黄", "萌", "萝", "菌", "菜", "萄", "菊", "萍", "菠", "营", "械", "梦", "梢", "梅", "检", "梳", "梯", "桶", "救", "副", "票", "戚", "爽", "聋", "袭", "盛", "雪", "辅", "辆", "虚", "雀", "堂", "常", "匙", "晨", "睁", "眯", "眼", "悬", "野", "啦", "晚", "啄", "距", "跃", "略", "蛇", "累", "唱", "患", "唯", "崖", "崭", "崇", "圈", "铜", "铲", "银", "甜", "梨", "犁", "移", "笨", "笼", "笛", "符", "第", "敏", "做", "袋", "悠", "偿", "偶", "偷", "您", "售", "停", "偏", "假", "得", "衔", "盘", "船", "斜", "盒", "鸽", "悉", "欲", "彩", "领", "脚", "脖", "脸", "脱", "象", "够", "猜", "猪", "猎", "猫", "猛", "馅", "馆", "凑", "减", "毫", "麻", "痒", "痕", "廊", "康", "庸", "鹿", "盗", "章", "竟", "商", "族", "旋", "望", "率", "着", "盖", "粘", "粗", "粒", "断", "剪", "兽", "清", "添", "淋", "淹", "渠", "渐", "混", "渔", "淘", "液", "淡", "深", "婆", "梁", "渗", "情", "惜", "惭", "悼", "惧", "惕", "惊", "惨", "惯", "寇", "寄", "宿", "窑", "密", "谋", "谎", "祸", "谜", "逮", "敢", "屠", "弹", "随", "蛋", "隆", "隐", "婚", "婶", "颈", "绩", "绪", "续", "骑", "绳", "维", "绵", "绸", "绿", "琴", "斑", "替", "款", "堪", "搭", "塔", "越", "趁", "趋", "超", "提", "堤", "博", "揭", "喜", "插", "揪", "搜", "煮", "援", "裁", "搁", "搂", "搅", "握", "揉", "斯", "期", "欺", "联", "散", "惹", "葬", "葛", "董", "葡", "敬", "葱", "落", "朝", "辜", "葵", "棒", "棋", "植", "森", "椅", "椒", "棵", "棍", "棉", "棚", "棕", "惠", "惑", "逼", "厨", "厦", "硬", "确", "雁", "殖", "裂", "雄", "暂", "雅", "辈", "悲", "紫", "辉", "敞", "赏", "掌", "晴", "暑", "最", "量", "喷", "晶", "喇", "遇", "喊", "景", "践", "跌", "跑", "遗", "蛙", "蛛", "蜓", "喝", "喂", "喘", "喉", "幅", "帽", "赌", "赔", "黑", "铸", "铺", "链", "销", "锁", "锄", "锅", "锈", "锋", "锐", "短", "智", "毯", "鹅", "剩", "稍", "程", "稀", "税", "筐", "等", "筑", "策", "筛", "筒", "答", "筋", "筝", "傲", "傅", "牌", "堡", "集", "焦", "傍", "储", "奥", "街", "惩", "御", "循", "艇", "舒", "番", "释", "禽", "腊", "脾", "腔", "鲁", "猾", "猴", "然", "馋", "装", "蛮", "就", "痛", "童", "阔", "善", "羡", "普", "粪", "尊", "道", "曾", "焰", "港", "湖", "渣", "湿", "温", "渴", "滑", "湾", "渡", "游", "滋", "溉", "愤", "慌", "惰", "愧", "愉", "慨", "割", "寒", "富", "窜", "窝", "窗", "遍", "裕", "裤", "裙", "谢", "谣", "谦", "属", "屡", "强", "粥", "疏", "隔", "隙", "絮", "嫂", "登", "缎", "缓", "编", "骗", "缘"}

type WordClickCaptchaGenerator struct {
	ResourceStore       resource.ImageCaptchaResourceStore
	ResourceImageReader *resource.ImageCaptchaResourceReaders
	ImageTransform      ImageTransform
	Fonts               []*truetype.Font
}

func (self *WordClickCaptchaGenerator) GenerateCaptchaImage(captchaExchange *model.CaptchaExchange) error {
	if self.Fonts == nil || len(self.Fonts) == 0 {
		return errors.Errorf("文字点选必须配置字体")
	}
	param := *captchaExchange.Param
	bgImageResource, err := resource.RequiredRandomGetResource(self.ResourceStore, param.CaptchaName, param.BackgroundImageTag)
	if err != nil {
		return errors.Wrap(err, "文字点选验证码获取背景图片失败")
	}
	bgImage, err := self.ResourceImageReader.GetResourceImage(bgImageResource)
	if err != nil {
		return errors.Wrap(err, "文字点选验证码转换背景图片失败")
	}
	checkFontCount := 4
	interferenceFontCount := 2
	fontImageWidth := bgImage.Bounds().Dx() / (checkFontCount + interferenceFontCount)
	fontSize := math.Max(70, float64(fontImageWidth)/math.Sqrt(2))
	clickImageArr := make([]*model.ClickImageCheckDefinition, 0, checkFontCount+interferenceFontCount)

	for i := 0; i < checkFontCount+interferenceFontCount; i++ {
		randomColor := common.GetRandomColor()
		// 随机汉字
		hanZi := chinese_chars[common.GetRandomInt(0, len(chinese_chars))]
		fontImage := common.NewTransparentImage(fontImageWidth, fontImageWidth)
		leftTop := (fontImageWidth - int(fontSize)) / 2
		fontImage = fontImage.DrawString(self.Fonts[common.GetRandomInt(0, len(self.Fonts))], randomColor, hanZi, leftTop, leftTop, fontSize)
		// 随机角度
		fontImage = fontImage.Rotate(float64(common.GetRandomInt(0, 360)))
		randomH := common.GetRandomInt(10, bgImage.Bounds().Dy()-fontImageWidth)
		var randomX int
		if i == 0 {
			randomX = 1
		} else {
			randomX = fontImageWidth * i
		}
		bgImage = bgImage.Overlay(fontImage, randomX, randomH)

		clickImageArr = append(clickImageArr, &model.ClickImageCheckDefinition{
			X:      randomX + fontImageWidth/2,
			Y:      randomH + fontImageWidth/2,
			Width:  fontImageWidth,
			Height: fontImageWidth,
			Value:  hanZi,
		})
	}

	shuffle(clickImageArr)
	clickImageArr = clickImageArr[:checkFontCount]
	// 生成点击图片
	var tipMsg string
	for _, c := range clickImageArr {
		tipMsg += c.Value
	}
	tipImage := common.NewTransparentImage(int(fontSize)*utf8.RuneCountInString(tipMsg), int(fontSize))
	tipImage = tipImage.DrawString(self.Fonts[common.GetRandomInt(0, len(self.Fonts))], color.RGBA{
		R: uint8(128),
		G: uint8(128),
		B: uint8(128),
		A: 255,
	}, tipMsg, 0, 0, fontSize)

	captchaExchange.BgImage = bgImage
	captchaExchange.ResourceImage = bgImageResource
	captchaExchange.TransferData = clickImageArr
	captchaExchange.TemplateImage = tipImage
	return nil
}

func (self *WordClickCaptchaGenerator) WrapImageCaptchaInfo(captchaExchange *model.CaptchaExchange) (*model.ImageCaptchaInfo, error) {
	clickImageArr := captchaExchange.TransferData.([]*model.ClickImageCheckDefinition)
	bgImage := captchaExchange.BgImage
	bgResource := captchaExchange.ResourceImage

	tipImage := captchaExchange.TemplateImage

	transFormData, err := self.ImageTransform.Transform(captchaExchange.Param, bgImage, tipImage, bgResource, nil, captchaExchange.CustomData)
	if err != nil {
		return nil, errors.Wrap(err, "图片转换异常")
	}
	customData := captchaExchange.CustomData
	customData.Expand = clickImageArr
	return &model.ImageCaptchaInfo{
		CaptchaName:           captchaExchange.Param.CaptchaName,
		CaptchaType:           model.TYPE_CLICK,
		BackgroundImage:       transFormData.BgImageUrl,
		TemplateImage:         transFormData.TemplateImageUrl,
		BackgroundImageTag:    bgResource.Tag,
		BackgroundImageHeight: bgImage.Bounds().Dy(),
		BackgroundImageWidth:  bgImage.Bounds().Dx(),
		TemplateImageWidth:    tipImage.Bounds().Dx(),
		TemplateImageHeight:   tipImage.Bounds().Dy(),
		Data:                  customData,
	}, nil
}

func shuffle(slice []*model.ClickImageCheckDefinition) {
	for i := range slice {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}