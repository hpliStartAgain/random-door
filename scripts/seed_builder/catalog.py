from __future__ import annotations

from copy import deepcopy


CITY_SLUGS = {
    "北京": "beijing",
    "西安": "xian",
    "南京": "nanjing",
    "杭州": "hangzhou",
    "成都": "chengdu",
    "广州": "guangzhou",
    "哈尔滨": "haerbin",
    "苏州": "suzhou",
    "大理": "dali",
    "厦门": "xiamen",
    "长沙": "changsha",
    "拉萨": "lasa",
    "上海": "shanghai",
    "重庆": "chongqing",
    "武汉": "wuhan",
    "天津": "tianjin",
    "青岛": "qingdao",
    "济南": "jinan",
    "郑州": "zhengzhou",
    "洛阳": "luoyang",
    "开封": "kaifeng",
    "合肥": "hefei",
    "福州": "fuzhou",
    "泉州": "quanzhou",
    "深圳": "shenzhen",
    "珠海": "zhuhai",
    "南宁": "nanning",
    "桂林": "guilin",
    "昆明": "kunming",
    "贵阳": "guiyang",
    "呼和浩特": "hohhot",
    "银川": "yinchuan",
    "兰州": "lanzhou",
    "乌鲁木齐": "urumqi",
    "三亚": "sanya",
}


EXTRA_CITIES = [
    {
        "name": "上海",
        "province": "上海",
        "lat": 31.2304,
        "lng": 121.4737,
        "intro": "海派文化窗口，黄浦江两岸连接近代风云与当代都市活力。",
        "dialect_sample": "侬好",
        "dialect_explanation": "上海话，意思是“你好”。",
        "tags": ["modern_city", "jiangnan", "coastal"],
        "landmarks": [
            {"name": "外滩", "description": "万国建筑群临江展开，是观察上海近代城市史的经典地标。"},
            {"name": "豫园", "description": "明代园林与城隍庙商圈相邻，呈现海派市井与江南园林气质。"},
        ],
        "foods": [
            {"name": "生煎馒头", "description": "底脆汤鲜，街头点心里最有上海日常烟火气。"},
            {"name": "本帮红烧肉", "description": "浓油赤酱，甜咸平衡，是本帮菜的代表味道。"},
        ],
        "character": {
            "name": "石库门讲述者",
            "character_type": "culture",
            "persona": "一位熟悉弄堂生活与海派文化的上海讲述者，说话利落细致。",
            "dialect_style": "上海话",
        },
    },
    {
        "name": "重庆",
        "province": "重庆",
        "lat": 29.5630,
        "lng": 106.5516,
        "intro": "山城沿江而起，码头文化、立体交通与火辣饮食交织成独特气质。",
        "dialect_sample": "要得",
        "dialect_explanation": "重庆话，表示可以、没问题。",
        "tags": ["spicy_food", "southwest", "modern_city"],
        "landmarks": [
            {"name": "洪崖洞", "description": "依山临江的吊脚楼风貌街区，夜景层次鲜明。"},
            {"name": "磁器口古镇", "description": "嘉陵江畔古镇，保留巴渝市井和码头记忆。"},
        ],
        "foods": [
            {"name": "重庆火锅", "description": "牛油红汤麻辣厚重，是山城餐桌的社交中心。"},
            {"name": "小面", "description": "佐料丰富、麻辣鲜香，是重庆人的日常早餐。"},
        ],
        "character": {
            "name": "山城摆谈人",
            "character_type": "culture",
            "persona": "一位熟悉重庆坡坎街巷的摆谈人，热情直接，爱讲山城生活。",
            "dialect_style": "重庆话",
        },
    },
    {
        "name": "武汉",
        "province": "湖北",
        "lat": 30.5928,
        "lng": 114.3055,
        "intro": "两江交汇、三镇并立，码头性格与高校文脉共同塑造江城气象。",
        "dialect_sample": "过早",
        "dialect_explanation": "武汉话，指吃早饭，也是一种热闹的城市生活方式。",
        "tags": ["river_city", "central_china", "modern_city"],
        "landmarks": [
            {"name": "黄鹤楼", "description": "长江边的文化名楼，承载诗文想象与城市地标意义。"},
            {"name": "东湖", "description": "城中大湖，绿道、水岸与高校片区相连。"},
        ],
        "foods": [
            {"name": "热干面", "description": "芝麻酱香浓，拌面爽滑，是武汉过早代表。"},
            {"name": "豆皮", "description": "糯米、蛋皮与馅料层叠，口感扎实香润。"},
        ],
        "character": {
            "name": "江城过早师傅",
            "character_type": "culture",
            "persona": "一位清晨就在街头忙碌的武汉过早师傅，爽快健谈。",
            "dialect_style": "武汉话",
        },
    },
    {
        "name": "天津",
        "province": "天津",
        "lat": 39.3434,
        "lng": 117.3616,
        "intro": "海河穿城而过，津门曲艺、租界建筑与码头商埠记忆并存。",
        "dialect_sample": "嘛玩意儿",
        "dialect_explanation": "天津话，表示“什么东西/怎么回事”，语气常带幽默感。",
        "tags": ["north_china", "coastal", "modern_city"],
        "landmarks": [
            {"name": "五大道", "description": "近代风貌建筑集中区，街区尺度适合慢慢步行观察。"},
            {"name": "天津之眼", "description": "跨河摩天轮，是海河夜景中的醒目标识。"},
        ],
        "foods": [
            {"name": "狗不理包子", "description": "传统津味名点，讲究褶花和鲜香馅料。"},
            {"name": "煎饼果子", "description": "绿豆面摊饼夹果篦，天津街头早餐代表。"},
        ],
        "character": {
            "name": "津门茶馆票友",
            "character_type": "culture",
            "persona": "一位爱听相声评书的天津茶馆票友，说话机灵幽默。",
            "dialect_style": "天津话",
        },
    },
    {
        "name": "青岛",
        "province": "山东",
        "lat": 36.0671,
        "lng": 120.3826,
        "intro": "红瓦绿树、碧海蓝天，德式街区、海港与啤酒文化共同成景。",
        "dialect_sample": "哈啤酒",
        "dialect_explanation": "青岛口语，意思是喝啤酒。",
        "tags": ["coastal", "north_china", "modern_city"],
        "landmarks": [
            {"name": "栈桥", "description": "伸入海湾的老地标，可远望红瓦城区与海岸线。"},
            {"name": "八大关", "description": "花木街巷与多国风格建筑构成青岛经典步行区域。"},
        ],
        "foods": [
            {"name": "鲅鱼水饺", "description": "海鲜馅料鲜嫩，是胶东沿海家常味道。"},
            {"name": "原浆啤酒", "description": "新鲜爽口，与海鲜排档一起构成城市体验。"},
        ],
        "character": {
            "name": "海边啤酒摊老板",
            "character_type": "culture",
            "persona": "一位熟悉海鲜市场和老街区的青岛摊主，爽朗实在。",
            "dialect_style": "青岛话",
        },
    },
    {
        "name": "济南",
        "province": "山东",
        "lat": 36.6512,
        "lng": 117.1201,
        "intro": "泉水塑造城市肌理，老城街巷、湖面与鲁菜香气相互映照。",
        "dialect_sample": "杠赛来",
        "dialect_explanation": "济南话，表示非常好、很厉害。",
        "tags": ["north_china", "spring_city", "food_city"],
        "landmarks": [
            {"name": "趵突泉", "description": "济南名泉代表，泉水涌动是城市最鲜明的视觉记忆。"},
            {"name": "大明湖", "description": "湖城相依的老城水景，连接泉水文化与市民生活。"},
        ],
        "foods": [
            {"name": "九转大肠", "description": "鲁菜经典，酸甜咸香层次复杂。"},
            {"name": "把子肉", "description": "酱香浓厚，常与米饭搭配，是济南日常快餐味道。"},
        ],
        "character": {
            "name": "泉城老茶客",
            "character_type": "culture",
            "persona": "一位常在泉边喝茶聊天的济南老茶客，朴实稳当。",
            "dialect_style": "济南话",
        },
    },
    {
        "name": "郑州",
        "province": "河南",
        "lat": 34.7466,
        "lng": 113.6254,
        "intro": "中原交通枢纽，商都遗址、黄河文化与现代城市建设并行。",
        "dialect_sample": "中不中",
        "dialect_explanation": "河南话，意思是行不行、可以吗。",
        "tags": ["central_china", "ancient_capital", "modern_city"],
        "landmarks": [
            {"name": "河南博物院", "description": "展示中原文明脉络，是理解河南历史的高密度入口。"},
            {"name": "二七纪念塔", "description": "郑州城市中心地标，记录近代铁路工人运动记忆。"},
        ],
        "foods": [
            {"name": "烩面", "description": "宽面筋道、羊汤浓郁，是河南面食代表。"},
            {"name": "胡辣汤", "description": "酸辣浓香，常与油馍头搭配做早餐。"},
        ],
        "character": {
            "name": "中原博物馆志愿者",
            "character_type": "culture",
            "persona": "一位熟悉中原文物和城市交通故事的郑州志愿者，表达清楚热心。",
            "dialect_style": "河南话",
        },
    },
    {
        "name": "洛阳",
        "province": "河南",
        "lat": 34.6197,
        "lng": 112.4540,
        "intro": "十三朝古都之一，牡丹、石窟与隋唐遗韵构成厚重城市底色。",
        "dialect_sample": "真得劲",
        "dialect_explanation": "河南话，表示很舒服、很有味道。",
        "tags": ["ancient_capital", "central_china"],
        "landmarks": [
            {"name": "龙门石窟", "description": "世界文化遗产，石刻造像展现北魏至唐代艺术高峰。"},
            {"name": "白马寺", "description": "中国佛教重要古寺，承载早期佛教传播记忆。"},
        ],
        "foods": [
            {"name": "洛阳水席", "description": "汤菜连台，酸辣鲜香，是洛阳宴席名片。"},
            {"name": "牛肉汤", "description": "清晨一碗热汤，是洛阳街头常见日常。"},
        ],
        "character": {
            "name": "洛阳牡丹花匠",
            "character_type": "culture",
            "persona": "一位懂牡丹也熟悉古都掌故的洛阳花匠，沉稳亲切。",
            "dialect_style": "河南话",
        },
    },
    {
        "name": "开封",
        "province": "河南",
        "lat": 34.7973,
        "lng": 114.3076,
        "intro": "北宋都城记忆鲜明，城摞城、夜市与汴梁风物延续古都烟火。",
        "dialect_sample": "可美",
        "dialect_explanation": "河南话，表示很好、很舒服。",
        "tags": ["ancient_capital", "central_china", "food_city"],
        "landmarks": [
            {"name": "清明上河园", "description": "以宋代市井图景为主题的园区，适合演示汴梁文化。"},
            {"name": "龙亭", "description": "开封老城轴线上的历史地标，可看湖景与古都轮廓。"},
        ],
        "foods": [
            {"name": "灌汤包", "description": "皮薄汤足，吃时讲究先开窗后喝汤。"},
            {"name": "桶子鸡", "description": "咸香紧实，是开封传统卤味代表。"},
        ],
        "character": {
            "name": "汴梁夜市掌灯人",
            "character_type": "culture",
            "persona": "一位熟悉开封夜市与宋文化的掌灯人，爱讲街巷故事。",
            "dialect_style": "河南话",
        },
    },
    {
        "name": "合肥",
        "province": "安徽",
        "lat": 31.8206,
        "lng": 117.2272,
        "intro": "江淮之间的科创城市，古逍遥津、巢湖风光与徽皖风味相连。",
        "dialect_sample": "搞得好",
        "dialect_explanation": "合肥口语，表示做得不错、挺好。",
        "tags": ["jianghuai", "modern_city", "food_city"],
        "landmarks": [
            {"name": "逍遥津公园", "description": "三国故事与城市公园结合，是合肥老城记忆之一。"},
            {"name": "巢湖", "description": "中国五大淡水湖之一，湖岸风光构成城市外延景观。"},
        ],
        "foods": [
            {"name": "李鸿章杂烩", "description": "徽菜名菜，食材丰富，汤味醇厚。"},
            {"name": "庐州烤鸭", "description": "皮香肉嫩，是合肥传统餐桌风味。"},
        ],
        "character": {
            "name": "庐州街巷讲解员",
            "character_type": "culture",
            "persona": "一位熟悉合肥旧称庐州和江淮饮食的讲解员，朴素耐心。",
            "dialect_style": "江淮官话",
        },
    },
    {
        "name": "福州",
        "province": "福建",
        "lat": 26.0745,
        "lng": 119.2965,
        "intro": "榕树成荫，三坊七巷、闽江水脉与茉莉花茶构成城市气息。",
        "dialect_sample": "虾油味",
        "dialect_explanation": "福州饮食口语，常形容带有本地调味风格的鲜咸味。",
        "tags": ["coastal", "minnan", "food_city"],
        "landmarks": [
            {"name": "三坊七巷", "description": "福州古城核心街区，保留坊巷格局与名人故居。"},
            {"name": "鼓山", "description": "临江名山，登高可看城市与闽江景致。"},
        ],
        "foods": [
            {"name": "佛跳墙", "description": "闽菜代表，汇聚海味与山珍，汤香浓郁。"},
            {"name": "鱼丸", "description": "外皮弹滑、内馅鲜香，是福州小吃名片。"},
        ],
        "character": {
            "name": "榕城茶铺老板",
            "character_type": "culture",
            "persona": "一位懂茉莉花茶和坊巷掌故的福州茶铺老板，温和细致。",
            "dialect_style": "福州话",
        },
    },
    {
        "name": "泉州",
        "province": "福建",
        "lat": 24.8741,
        "lng": 118.6759,
        "intro": "宋元海丝重要港口，多元宗教遗存、闽南红砖厝与南音相映。",
        "dialect_sample": "爱拼才会赢",
        "dialect_explanation": "闽南语流行表达，意思是努力拼搏才有机会成功。",
        "tags": ["coastal", "minnan", "heritage"],
        "landmarks": [
            {"name": "开元寺", "description": "泉州古刹，东西塔是城市天际线的重要标志。"},
            {"name": "洛阳桥", "description": "中国古代海港桥梁代表，体现宋代工程智慧。"},
        ],
        "foods": [
            {"name": "面线糊", "description": "细软顺滑，可加海蛎、大肠等配料。"},
            {"name": "姜母鸭", "description": "姜香浓郁，砂锅慢煲，是闽南经典菜。"},
        ],
        "character": {
            "name": "海丝南音艺人",
            "character_type": "culture",
            "persona": "一位会唱南音、熟悉海丝遗存的泉州艺人，谦和从容。",
            "dialect_style": "闽南语",
        },
    },
    {
        "name": "深圳",
        "province": "广东",
        "lat": 22.5431,
        "lng": 114.0579,
        "intro": "改革开放窗口，年轻移民城市在海岸线、高楼与创新产业中生长。",
        "dialect_sample": "搞掂",
        "dialect_explanation": "粤语常用语，表示搞定、解决了。",
        "tags": ["modern_city", "coastal", "lingnan"],
        "landmarks": [
            {"name": "莲花山公园", "description": "城市中心绿地，可俯瞰福田天际线。"},
            {"name": "大鹏所城", "description": "明清海防古城，展示深圳更早的海防记忆。"},
        ],
        "foods": [
            {"name": "早茶点心", "description": "深圳融合广府饮食，早茶是重要社交场景。"},
            {"name": "光明乳鸽", "description": "皮脆肉嫩，是深圳本地特色菜之一。"},
        ],
        "character": {
            "name": "南山产品经理",
            "character_type": "symbol",
            "persona": "一位在深圳工作的年轻产品经理，熟悉创新产业和城市生活。",
            "dialect_style": "普通话夹少量粤语",
        },
    },
    {
        "name": "珠海",
        "province": "广东",
        "lat": 22.2711,
        "lng": 113.5767,
        "intro": "滨海花园城市，情侣路、海岛与粤澳相邻的口岸气质鲜明。",
        "dialect_sample": "靓",
        "dialect_explanation": "粤语，表示漂亮、好看。",
        "tags": ["coastal", "lingnan", "modern_city"],
        "landmarks": [
            {"name": "情侣路", "description": "沿海步道连接城市海景，是珠海最具识别度的公共空间。"},
            {"name": "珠海渔女", "description": "海边雕塑地标，象征城市滨海形象。"},
        ],
        "foods": [
            {"name": "横琴蚝", "description": "肉质肥美，体现珠海海鲜风味。"},
            {"name": "斗门重壳蟹", "description": "珠海本地水产名物，鲜甜浓郁。"},
        ],
        "character": {
            "name": "情侣路骑行者",
            "character_type": "culture",
            "persona": "一位熟悉海岸线和海岛玩法的珠海骑行者，轻松友好。",
            "dialect_style": "粤语",
        },
    },
    {
        "name": "南宁",
        "province": "广西",
        "lat": 22.8170,
        "lng": 108.3669,
        "intro": "绿城南宁连接壮乡风情与东盟交流，骑楼、夜市和米粉香气浓郁。",
        "dialect_sample": "靓仔靓女",
        "dialect_explanation": "粤桂地区常用称呼，表示帅哥美女，也可作亲切称呼。",
        "tags": ["lingnan", "southwest", "food_city"],
        "landmarks": [
            {"name": "青秀山", "description": "城市绿肺，山水园林和观景平台适合慢游。"},
            {"name": "三街两巷", "description": "老城更新街区，展示骑楼与南宁历史生活。"},
        ],
        "foods": [
            {"name": "老友粉", "description": "酸辣开胃，蒜香豆豉味突出，是南宁代表米粉。"},
            {"name": "柠檬鸭", "description": "酸香清爽，体现广西菜的明亮风味。"},
        ],
        "character": {
            "name": "绿城米粉摊主",
            "character_type": "culture",
            "persona": "一位在夜市卖米粉的南宁摊主，热情爽快，熟悉壮乡风味。",
            "dialect_style": "南宁白话",
        },
    },
    {
        "name": "桂林",
        "province": "广西",
        "lat": 25.2736,
        "lng": 110.2900,
        "intro": "喀斯特山水名城，漓江、象鼻山与米粉构成高度可感的旅行记忆。",
        "dialect_sample": "蛮好耍",
        "dialect_explanation": "桂林口语，表示很好玩。",
        "tags": ["southwest", "landscape", "food_city"],
        "landmarks": [
            {"name": "漓江", "description": "山水画卷般的河流景观，是桂林旅游核心名片。"},
            {"name": "象鼻山", "description": "形似巨象饮水，是桂林城市地标。"},
        ],
        "foods": [
            {"name": "桂林米粉", "description": "卤水香浓，配锅烧、酸豆角等小料。"},
            {"name": "啤酒鱼", "description": "阳朔风味菜，鱼肉鲜嫩、酱香微辣。"},
        ],
        "character": {
            "name": "漓江竹筏船工",
            "character_type": "culture",
            "persona": "一位熟悉漓江山水和村镇故事的竹筏船工，朴实健谈。",
            "dialect_style": "桂林话",
        },
    },
    {
        "name": "昆明",
        "province": "云南",
        "lat": 24.8801,
        "lng": 102.8329,
        "intro": "春城气候温和，滇池、西山、花市与多民族饮食构成明亮城市感。",
        "dialect_sample": "板扎",
        "dialect_explanation": "云南话，表示很好、很棒。",
        "tags": ["southwest", "modern_city", "food_city"],
        "landmarks": [
            {"name": "滇池", "description": "高原湖泊，海埂大坝和西山共同构成昆明风景线。"},
            {"name": "金马碧鸡坊", "description": "老城地标，连接昆明商业街区与历史记忆。"},
        ],
        "foods": [
            {"name": "过桥米线", "description": "热汤烫熟配料，汤鲜料足，是云南代表小吃。"},
            {"name": "鲜花饼", "description": "玫瑰花馅清香，是昆明伴手礼代表。"},
        ],
        "character": {
            "name": "春城花市老板",
            "character_type": "culture",
            "persona": "一位熟悉鲜花、米线和高原湖景的昆明花市老板，开朗温柔。",
            "dialect_style": "云南话",
        },
    },
    {
        "name": "贵阳",
        "province": "贵州",
        "lat": 26.6470,
        "lng": 106.6302,
        "intro": "山地城市清爽宜居，酸辣饮食、民族文化和夜市烟火气突出。",
        "dialect_sample": "安逸",
        "dialect_explanation": "西南地区常用语，表示舒服、惬意。",
        "tags": ["southwest", "spicy_food", "food_city"],
        "landmarks": [
            {"name": "甲秀楼", "description": "南明河畔古楼，是贵阳老城最具代表性的地标。"},
            {"name": "黔灵山公园", "description": "城市山地公园，寺庙、湖水与山林相连。"},
        ],
        "foods": [
            {"name": "丝娃娃", "description": "薄饼卷多样蔬菜，蘸酸辣汤汁，清爽开胃。"},
            {"name": "酸汤鱼", "description": "酸香鲜辣，体现贵州饮食的发酵风味。"},
        ],
        "character": {
            "name": "黔城酸汤厨娘",
            "character_type": "culture",
            "persona": "一位熟悉贵州酸辣风味和夜市生活的厨娘，热情细心。",
            "dialect_style": "贵阳话",
        },
    },
    {
        "name": "呼和浩特",
        "province": "内蒙古",
        "lat": 40.8426,
        "lng": 111.7492,
        "intro": "草原文化与塞北城市生活交汇，寺庙、乳茶和蒙古族风情鲜明。",
        "dialect_sample": "赛白努",
        "dialect_explanation": "蒙古语问候语，意思是你好。",
        "tags": ["north_china", "grassland", "food_city"],
        "landmarks": [
            {"name": "大召寺", "description": "呼和浩特重要寺庙，体现藏传佛教与城市历史。"},
            {"name": "内蒙古博物院", "description": "展示草原文明、自然生态与民族文化。"},
        ],
        "foods": [
            {"name": "手把肉", "description": "羊肉原味鲜香，是草原饮食代表。"},
            {"name": "奶茶", "description": "咸香温热，常与奶食和肉食搭配。"},
        ],
        "character": {
            "name": "草原博物馆讲解员",
            "character_type": "culture",
            "persona": "一位熟悉草原文化和呼和浩特城市历史的讲解员，稳重亲切。",
            "dialect_style": "普通话夹少量蒙古语问候",
        },
    },
    {
        "name": "银川",
        "province": "宁夏",
        "lat": 38.4872,
        "lng": 106.2309,
        "intro": "贺兰山下的塞上湖城，西夏遗存、黄河灌溉与清真美食相连。",
        "dialect_sample": "攒劲",
        "dialect_explanation": "西北口语，表示很棒、有劲道。",
        "tags": ["northwest", "food_city", "heritage"],
        "landmarks": [
            {"name": "西夏陵", "description": "西夏王陵遗址，展现西北历史文化的重要面向。"},
            {"name": "镇北堡西部影城", "description": "电影取景地与西北风貌结合，适合演示文化记忆。"},
        ],
        "foods": [
            {"name": "羊杂碎", "description": "汤鲜味重，是宁夏街头常见早餐。"},
            {"name": "手抓羊肉", "description": "肉质鲜嫩，突出西北羊肉原香。"},
        ],
        "character": {
            "name": "贺兰山脚驼队向导",
            "character_type": "culture",
            "persona": "一位熟悉贺兰山、西夏遗存和宁夏饮食的向导，沉着可靠。",
            "dialect_style": "宁夏话",
        },
    },
    {
        "name": "兰州",
        "province": "甘肃",
        "lat": 36.0611,
        "lng": 103.8343,
        "intro": "黄河穿城而过，铁桥、牛肉面和西北通道气质塑造金城印象。",
        "dialect_sample": "尕",
        "dialect_explanation": "兰州及西北口语，常表示小、亲切。",
        "tags": ["northwest", "food_city", "river_city"],
        "landmarks": [
            {"name": "中山桥", "description": "黄河铁桥，是兰州最醒目的城市地标之一。"},
            {"name": "白塔山", "description": "隔河俯瞰城区与黄河，是理解兰州地形的好位置。"},
        ],
        "foods": [
            {"name": "兰州牛肉面", "description": "一清二白三红四绿，汤清面劲道。"},
            {"name": "灰豆子", "description": "甜香绵密，是兰州传统小吃。"},
        ],
        "character": {
            "name": "黄河边拉面师傅",
            "character_type": "culture",
            "persona": "一位清晨拉面的兰州师傅，熟悉黄河边生活和西北饮食。",
            "dialect_style": "兰州话",
        },
    },
    {
        "name": "乌鲁木齐",
        "province": "新疆",
        "lat": 43.8256,
        "lng": 87.6168,
        "intro": "天山脚下的多民族城市，大巴扎、烤肉香气与西域通道记忆鲜明。",
        "dialect_sample": "亚克西",
        "dialect_explanation": "维吾尔语常用表达，意思是好、不错。",
        "tags": ["northwest", "silk_road", "food_city"],
        "landmarks": [
            {"name": "国际大巴扎", "description": "集市、建筑和演艺结合，是观察新疆多元文化的窗口。"},
            {"name": "红山公园", "description": "城市中心山体公园，可俯瞰乌鲁木齐城区。"},
        ],
        "foods": [
            {"name": "烤羊肉串", "description": "炭火香气浓郁，孜然与辣椒突出。"},
            {"name": "大盘鸡", "description": "鸡肉、土豆和皮带面组合，分量扎实。"},
        ],
        "character": {
            "name": "天山巴扎商人",
            "character_type": "culture",
            "persona": "一位熟悉大巴扎和丝路风物的商人，热情大方。",
            "dialect_style": "普通话夹少量维吾尔语问候",
        },
    },
    {
        "name": "三亚",
        "province": "海南",
        "lat": 18.2528,
        "lng": 109.5119,
        "intro": "热带滨海度假城市，海湾、椰风、黎苗文化和海鲜夜市构成南端体验。",
        "dialect_sample": "椰香",
        "dialect_explanation": "海南旅行语境中常用来形容带有椰子清香的风味。",
        "tags": ["coastal", "lingnan", "landscape"],
        "landmarks": [
            {"name": "亚龙湾", "description": "海水清澈、沙滩绵长，是三亚代表性海湾。"},
            {"name": "天涯海角", "description": "海滨礁石景区，承载浪漫旅行想象。"},
        ],
        "foods": [
            {"name": "海南鸡饭", "description": "鸡肉嫩滑、米饭有鸡油香，是海南代表味道。"},
            {"name": "清补凉", "description": "椰奶、豆类和水果组合，清爽解暑。"},
        ],
        "character": {
            "name": "海湾冲浪教练",
            "character_type": "culture",
            "persona": "一位熟悉三亚海湾和本地夜市的冲浪教练，阳光直接。",
            "dialect_style": "海南话问候夹普通话",
        },
    },
]


def normalize_existing_city(city: dict) -> dict:
    """Preserve existing copy while normalizing image paths to generated PNG assets."""
    normalized = deepcopy(city)
    slug = CITY_SLUGS[city["name"]]
    normalized["cover_image_url"] = f"/static/landmarks/{slug}_cover.png"
    for index, landmark in enumerate(normalized.get("landmarks", []), start=1):
        landmark["image_url"] = f"/static/landmarks/{slug}_landmark_{index}.png"
    for index, food in enumerate(normalized.get("foods", []), start=1):
        food["image_url"] = f"/static/foods/{slug}_food_{index}.png"
    for character in normalized.get("characters", []):
        character["avatar_url"] = f"/static/characters/{slug}_character.png"
    return normalized


def build_extra_city(source: dict) -> dict:
    slug = CITY_SLUGS[source["name"]]
    city = {
        "name": source["name"],
        "province": source["province"],
        "lat": source["lat"],
        "lng": source["lng"],
        "intro": source["intro"],
        "cover_image_url": f"/static/landmarks/{slug}_cover.png",
        "dialect_sample": source["dialect_sample"],
        "dialect_explanation": source["dialect_explanation"],
        "tags": source["tags"],
        "landmarks": [],
        "foods": [],
        "characters": [],
    }
    for index, landmark in enumerate(source["landmarks"], start=1):
        city["landmarks"].append(
            {
                "name": landmark["name"],
                "image_url": f"/static/landmarks/{slug}_landmark_{index}.png",
                "description": landmark["description"],
            }
        )
    for index, food in enumerate(source["foods"], start=1):
        city["foods"].append(
            {
                "name": food["name"],
                "image_url": f"/static/foods/{slug}_food_{index}.png",
                "description": food["description"],
            }
        )
    character = source["character"]
    city["characters"].append(
        {
            "name": character["name"],
            "character_type": character["character_type"],
            "avatar_url": f"/static/characters/{slug}_character.png",
            "persona": character["persona"],
            "dialect_style": character["dialect_style"],
            "prompt": (
                f"你是{character['name']}，基于{source['name']}城市文化、地标与美食进行友好对话；"
                "可用少量方言并附普通话解释，不编史、不声称真实复活，回答150字内。"
            ),
        }
    )
    return city


def build_catalog(existing_cities: list[dict]) -> list[dict]:
    preserved = [normalize_existing_city(city) for city in existing_cities]
    existing_names = {city["name"] for city in preserved}
    extras = [build_extra_city(city) for city in EXTRA_CITIES if city["name"] not in existing_names]
    return preserved + extras
