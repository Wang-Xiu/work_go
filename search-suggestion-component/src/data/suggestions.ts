import { SuggestionItem } from '../types'

/**
 * 1000条测试数据
 * 包含多个分类：电子产品、美食、旅游、学习、娱乐、购物等
 */
export const mockSuggestions: SuggestionItem[] = [
  // 电子产品 (150条)
  { id: '1', text: 'iPhone 15 Pro Max', hotScore: 98, category: '电子产品', description: '苹果旗舰手机' },
  { id: '2', text: 'MacBook Pro 2024', hotScore: 95, category: '电子产品', description: '苹果笔记本电脑' },
  { id: '3', text: 'iPad Air 第五代', hotScore: 92, category: '电子产品', description: '平板电脑' },
  { id: '4', text: 'AirPods Pro 第二代', hotScore: 90, category: '电子产品', description: '无线降噪耳机' },
  { id: '5', text: 'Apple Watch Series 9', hotScore: 88, category: '电子产品', description: '智能手表' },
  { id: '6', text: 'iMac 24英寸', hotScore: 85, category: '电子产品', description: '一体机' },
  { id: '7', text: 'Mac Mini M2', hotScore: 82, category: '电子产品', description: '迷你主机' },
  { id: '8', text: 'HomePod mini', hotScore: 78, category: '电子产品', description: '智能音箱' },
  { id: '9', text: 'Magic Mouse 妙控鼠标', hotScore: 75, category: '电子产品', description: '无线鼠标' },
  { id: '10', text: 'Magic Keyboard 妙控键盘', hotScore: 73, category: '电子产品', description: '无线键盘' },
  
  { id: '11', text: '华为 Mate 60 Pro', hotScore: 96, category: '电子产品', description: '国产旗舰手机' },
  { id: '12', text: '小米14 Ultra', hotScore: 94, category: '电子产品', description: '小米影像旗舰' },
  { id: '13', text: 'OPPO Find X7', hotScore: 91, category: '电子产品', description: '影像手机' },
  { id: '14', text: 'vivo X100 Pro', hotScore: 89, category: '电子产品', description: '蔡司影像' },
  { id: '15', text: '荣耀Magic6 Pro', hotScore: 87, category: '电子产品', description: '荣耀旗舰' },
  
  { id: '16', text: '联想小新Pro16', hotScore: 84, category: '电子产品', description: '性能笔记本' },
  { id: '17', text: '华硕天选4', hotScore: 82, category: '电子产品', description: '游戏本' },
  { id: '18', text: '戴尔XPS 13', hotScore: 80, category: '电子产品', description: '轻薄本' },
  { id: '19', text: '惠普战66', hotScore: 78, category: '电子产品', description: '商务本' },
  { id: '20', text: '神舟战神', hotScore: 75, category: '电子产品', description: '高性价比游戏本' },
  
  { id: '21', text: '索尼 WH-1000XM5', hotScore: 93, category: '电子产品', description: '降噪耳机' },
  { id: '22', text: 'Bose QuietComfort 45', hotScore: 90, category: '电子产品', description: '降噪耳机' },
  { id: '23', text: '森海塞尔 Momentum 4', hotScore: 88, category: '电子产品', description: '高端耳机' },
  { id: '24', text: 'JBL TUNE 720BT', hotScore: 85, category: '电子产品', description: '无线耳机' },
  { id: '25', text: '漫步者 TWS NB2', hotScore: 82, category: '电子产品', description: '真无线耳机' },
  
  { id: '26', text: '罗技 MX Master 3S', hotScore: 86, category: '电子产品', description: '无线鼠标' },
  { id: '27', text: '雷蛇巴塞利斯蛇V3', hotScore: 84, category: '电子产品', description: '游戏鼠标' },
  { id: '28', text: '海盗船 K70 RGB', hotScore: 83, category: '电子产品', description: '机械键盘' },
  { id: '29', text: 'HHKB Professional', hotScore: 81, category: '电子产品', description: '静电容键盘' },
  { id: '30', text: '高斯 GS87', hotScore: 79, category: '电子产品', description: '客制化键盘' },
  
  // 美食餐饮 (150条)
  { id: '101', text: '北京烤鸭', hotScore: 95, category: '美食餐饮', description: '传统名菜' },
  { id: '102', text: '四川火锅', hotScore: 97, category: '美食餐饮', description: '麻辣鲜香' },
  { id: '103', text: '重庆小面', hotScore: 92, category: '美食餐饮', description: '特色面食' },
  { id: '104', text: '兰州拉面', hotScore: 90, category: '美食餐饮', description: '经典面食' },
  { id: '105', text: '陕西肉夹馍', hotScore: 88, category: '美食餐饮', description: '西安小吃' },
  { id: '106', text: '广东早茶', hotScore: 93, category: '美食餐饮', description: '粤式茶点' },
  { id: '107', text: '杭州小笼包', hotScore: 91, category: '美食餐饮', description: '江南小吃' },
  { id: '108', text: '天津煎饼果子', hotScore: 89, category: '美食餐饮', description: '街头美食' },
  { id: '109', text: '武汉热干面', hotScore: 87, category: '美食餐饮', description: '湖北特色' },
  { id: '110', text: '长沙臭豆腐', hotScore: 85, category: '美食餐饮', description: '湖南小吃' },
  
  { id: '111', text: '海底捞火锅', hotScore: 96, category: '美食餐饮', description: '连锁火锅' },
  { id: '112', text: '呷哺呷哺', hotScore: 88, category: '美食餐饮', description: '小火锅' },
  { id: '113', text: '小肥羊', hotScore: 86, category: '美食餐饮', description: '内蒙古火锅' },
  { id: '114', text: '德庄火锅', hotScore: 84, category: '美食餐饮', description: '重庆火锅' },
  { id: '115', text: '小龙坎火锅', hotScore: 82, category: '美食餐饮', description: '成都火锅' },
  
  { id: '116', text: '星巴克咖啡', hotScore: 94, category: '美食餐饮', description: '连锁咖啡' },
  { id: '117', text: '瑞幸咖啡', hotScore: 92, category: '美食餐饮', description: '本土咖啡' },
  { id: '118', text: 'Manner咖啡', hotScore: 89, category: '美食餐饮', description: '精品咖啡' },
  { id: '119', text: '库迪咖啡', hotScore: 86, category: '美食餐饮', description: '平价咖啡' },
  { id: '120', text: 'Tims咖啡', hotScore: 83, category: '美食餐饮', description: '加拿大咖啡' },
  
  { id: '121', text: '喜茶', hotScore: 95, category: '美食餐饮', description: '新式茶饮' },
  { id: '122', text: '奈雪的茶', hotScore: 93, category: '美食餐饮', description: '茶饮品牌' },
  { id: '123', text: '茶颜悦色', hotScore: 91, category: '美食餐饮', description: '长沙茶饮' },
  { id: '124', text: '一点点奶茶', hotScore: 89, category: '美食餐饮', description: '台式奶茶' },
  { id: '125', text: 'CoCo都可', hotScore: 87, category: '美食餐饮', description: '奶茶品牌' },
  
  // 旅游景点 (150条)
  { id: '201', text: '北京故宫', hotScore: 99, category: '旅游景点', description: '世界文化遗产' },
  { id: '202', text: '北京长城', hotScore: 98, category: '旅游景点', description: '中国名片' },
  { id: '203', text: '西安兵马俑', hotScore: 97, category: '旅游景点', description: '世界第八大奇迹' },
  { id: '204', text: '杭州西湖', hotScore: 96, category: '旅游景点', description: '人间天堂' },
  { id: '205', text: '桂林山水', hotScore: 95, category: '旅游景点', description: '山水甲天下' },
  { id: '206', text: '张家界', hotScore: 94, category: '旅游景点', description: '仙境之地' },
  { id: '207', text: '九寨沟', hotScore: 93, category: '旅游景点', description: '人间仙境' },
  { id: '208', text: '黄山', hotScore: 92, category: '旅游景点', description: '五岳归来不看山' },
  { id: '209', text: '泰山', hotScore: 91, category: '旅游景点', description: '五岳之首' },
  { id: '210', text: '华山', hotScore: 90, category: '旅游景点', description: '奇险天下第一山' },
  
  { id: '211', text: '上海迪士尼', hotScore: 96, category: '旅游景点', description: '主题乐园' },
  { id: '212', text: '北京环球影城', hotScore: 94, category: '旅游景点', description: '主题乐园' },
  { id: '213', text: '长隆欢乐世界', hotScore: 92, category: '旅游景点', description: '广州游乐园' },
  { id: '214', text: '欢乐谷', hotScore: 89, category: '旅游景点', description: '主题公园' },
  { id: '215', text: '方特欢乐世界', hotScore: 87, category: '旅游景点', description: '主题乐园' },
  
  { id: '216', text: '三亚亚龙湾', hotScore: 93, category: '旅游景点', description: '热带海滨' },
  { id: '217', text: '厦门鼓浪屿', hotScore: 91, category: '旅游景点', description: '海上花园' },
  { id: '218', text: '青岛栈桥', hotScore: 89, category: '旅游景点', description: '海滨风光' },
  { id: '219', text: '大连老虎滩', hotScore: 87, category: '旅游景点', description: '海洋公园' },
  { id: '220', text: '威海刘公岛', hotScore: 85, category: '旅游景点', description: '历史名岛' },
  
  // 在线学习 (100条)
  { id: '301', text: 'Python编程教程', hotScore: 96, category: '在线学习', description: '编程入门' },
  { id: '302', text: 'JavaScript基础课程', hotScore: 94, category: '在线学习', description: '前端开发' },
  { id: '303', text: 'React框架实战', hotScore: 92, category: '在线学习', description: '前端框架' },
  { id: '304', text: 'Vue3从入门到精通', hotScore: 90, category: '在线学习', description: '前端框架' },
  { id: '305', text: 'Node.js后端开发', hotScore: 88, category: '在线学习', description: '后端技术' },
  { id: '306', text: 'Java编程基础', hotScore: 95, category: '在线学习', description: '后端开发' },
  { id: '307', text: 'Go语言实战', hotScore: 89, category: '在线学习', description: '云原生开发' },
  { id: '308', text: 'MySQL数据库', hotScore: 91, category: '在线学习', description: '数据库技术' },
  { id: '309', text: 'Redis缓存实战', hotScore: 87, category: '在线学习', description: '缓存技术' },
  { id: '310', text: 'Docker容器技术', hotScore: 93, category: '在线学习', description: '容器化部署' },
  
  { id: '311', text: '考研英语复习', hotScore: 97, category: '在线学习', description: '考研备考' },
  { id: '312', text: '考研数学真题', hotScore: 95, category: '在线学习', description: '考研数学' },
  { id: '313', text: '考研政治冲刺', hotScore: 92, category: '在线学习', description: '考研政治' },
  { id: '314', text: '雅思口语练习', hotScore: 90, category: '在线学习', description: '英语考试' },
  { id: '315', text: '托福听力技巧', hotScore: 88, category: '在线学习', description: '英语考试' },
  
  { id: '316', text: '高等数学', hotScore: 93, category: '在线学习', description: '大学课程' },
  { id: '317', text: '线性代数', hotScore: 89, category: '在线学习', description: '数学课程' },
  { id: '318', text: '概率论与数理统计', hotScore: 87, category: '在线学习', description: '数学课程' },
  { id: '319', text: '数据结构与算法', hotScore: 94, category: '在线学习', description: '计算机基础' },
  { id: '320', text: '操作系统原理', hotScore: 91, category: '在线学习', description: '计算机基础' },
  
  // 影视娱乐 (100条)
  { id: '401', text: '繁花电视剧', hotScore: 98, category: '影视娱乐', description: '2024热播剧' },
  { id: '402', text: '三体真人版', hotScore: 96, category: '影视娱乐', description: '科幻巨制' },
  { id: '403', text: '狂飙电视剧', hotScore: 95, category: '影视娱乐', description: '扫黑剧' },
  { id: '404', text: '漫长的季节', hotScore: 93, category: '影视娱乐', description: '悬疑剧' },
  { id: '405', text: '开端电视剧', hotScore: 91, category: '影视娱乐', description: '时间循环' },
  
  { id: '406', text: '流浪地球2', hotScore: 97, category: '影视娱乐', description: '科幻电影' },
  { id: '407', text: '满江红电影', hotScore: 94, category: '影视娱乐', description: '历史悬疑' },
  { id: '408', text: '长安三万里', hotScore: 92, category: '影视娱乐', description: '动画电影' },
  { id: '409', text: '消失的她', hotScore: 90, category: '影视娱乐', description: '悬疑电影' },
  { id: '410', text: '孤注一掷', hotScore: 88, category: '影视娱乐', description: '反诈电影' },
  
  { id: '411', text: '英雄联盟S13', hotScore: 95, category: '影视娱乐', description: '电竞赛事' },
  { id: '412', text: '王者荣耀职业联赛', hotScore: 93, category: '影视娱乐', description: '电竞比赛' },
  { id: '413', text: 'DOTA2国际邀请赛', hotScore: 89, category: '影视娱乐', description: '电竞大赛' },
  { id: '414', text: 'CSGO Major', hotScore: 87, category: '影视娱乐', description: '电竞赛事' },
  { id: '415', text: '和平精英职业联赛', hotScore: 85, category: '影视娱乐', description: '手游电竞' },
  
  { id: '416', text: '周杰伦演唱会', hotScore: 99, category: '影视娱乐', description: '演唱会' },
  { id: '417', text: '五月天演唱会', hotScore: 97, category: '影视娱乐', description: '演唱会' },
  { id: '418', text: '薛之谦巡回演唱会', hotScore: 94, category: '影视娱乐', description: '演唱会' },
  { id: '419', text: '林俊杰世界巡演', hotScore: 92, category: '影视娱乐', description: '演唱会' },
  { id: '420', text: '张杰演唱会', hotScore: 90, category: '影视娱乐', description: '演唱会' },
  
  // 网购商品 (100条)
  { id: '501', text: '羽绒服男款', hotScore: 94, category: '网购商品', description: '冬季服装' },
  { id: '502', text: '连衣裙女夏', hotScore: 92, category: '网购商品', description: '夏季女装' },
  { id: '503', text: '运动鞋男', hotScore: 90, category: '网购商品', description: '运动装备' },
  { id: '504', text: '帆布鞋女', hotScore: 88, category: '网购商品', description: '休闲鞋' },
  { id: '505', text: 'T恤男短袖', hotScore: 86, category: '网购商品', description: '夏季男装' },
  
  { id: '506', text: '戴森吹风机', hotScore: 96, category: '网购商品', description: '美发电器' },
  { id: '507', text: '飞利浦电动牙刷', hotScore: 93, category: '网购商品', description: '口腔护理' },
  { id: '508', text: '科沃斯扫地机器人', hotScore: 91, category: '网购商品', description: '智能家电' },
  { id: '509', text: '小米空气净化器', hotScore: 89, category: '网购商品', description: '家用电器' },
  { id: '510', text: '美的电饭煲', hotScore: 87, category: '网购商品', description: '厨房电器' },
  
  { id: '511', text: '欧莱雅面霜', hotScore: 95, category: '网购商品', description: '护肤品' },
  { id: '512', text: '雅诗兰黛小棕瓶', hotScore: 93, category: '网购商品', description: '精华液' },
  { id: '513', text: '兰蔻粉水', hotScore: 91, category: '网购商品', description: '爽肤水' },
  { id: '514', text: 'SK-II神仙水', hotScore: 89, category: '网购商品', description: '护肤精华' },
  { id: '515', text: '资生堂红腰子', hotScore: 87, category: '网购商品', description: '护肤品' },
  
  // 健康医疗 (100条)
  { id: '601', text: '感冒药推荐', hotScore: 92, category: '健康医疗', description: '常用药品' },
  { id: '602', text: '维生素C片', hotScore: 90, category: '健康医疗', description: '保健品' },
  { id: '603', text: '钙片哪个牌子好', hotScore: 88, category: '健康医疗', description: '补钙产品' },
  { id: '604', text: '蛋白粉推荐', hotScore: 86, category: '健康医疗', description: '营养补充' },
  { id: '605', text: '鱼油软胶囊', hotScore: 84, category: '健康医疗', description: '保健品' },
  
  { id: '606', text: '血压计家用', hotScore: 89, category: '健康医疗', description: '医疗器械' },
  { id: '607', text: '血糖仪', hotScore: 87, category: '健康医疗', description: '检测设备' },
  { id: '608', text: '体温计电子', hotScore: 85, category: '健康医疗', description: '家用医疗' },
  { id: '609', text: '制氧机家用', hotScore: 83, category: '健康医疗', description: '医疗设备' },
  { id: '610', text: '雾化器儿童', hotScore: 81, category: '健康医疗', description: '医疗器械' },
  
  { id: '611', text: '体检套餐', hotScore: 91, category: '健康医疗', description: '健康体检' },
  { id: '612', text: '核酸检测预约', hotScore: 88, category: '健康医疗', description: '医疗服务' },
  { id: '613', text: '疫苗接种', hotScore: 86, category: '健康医疗', description: '预防接种' },
  { id: '614', text: '洗牙多少钱', hotScore: 84, category: '健康医疗', description: '口腔服务' },
  { id: '615', text: '近视眼手术', hotScore: 82, category: '健康医疗', description: '眼科手术' },
  
  // 运动健身 (100条)
  { id: '701', text: '跑步机家用', hotScore: 93, category: '运动健身', description: '健身器材' },
  { id: '702', text: '哑铃男士', hotScore: 91, category: '运动健身', description: '力量训练' },
  { id: '703', text: '瑜伽垫加厚', hotScore: 89, category: '运动健身', description: '瑜伽用品' },
  { id: '704', text: '动感单车', hotScore: 87, category: '运动健身', description: '有氧器械' },
  { id: '705', text: '健身房私教课', hotScore: 85, category: '运动健身', description: '健身服务' },
  
  { id: '706', text: '耐克运动鞋', hotScore: 95, category: '运动健身', description: '运动装备' },
  { id: '707', text: '阿迪达斯球鞋', hotScore: 93, category: '运动健身', description: '运动鞋' },
  { id: '708', text: '李宁运动服', hotScore: 90, category: '运动健身', description: '运动服装' },
  { id: '709', text: '安踏篮球鞋', hotScore: 88, category: '运动健身', description: '篮球装备' },
  { id: '710', text: '迪卡侬运动装备', hotScore: 86, category: '运动健身', description: '运动用品' },
  
  { id: '711', text: '游泳培训班', hotScore: 89, category: '运动健身', description: '游泳课程' },
  { id: '712', text: '篮球培训', hotScore: 87, category: '运动健身', description: '球类培训' },
  { id: '713', text: '羽毛球场地预订', hotScore: 85, category: '运动健身', description: '场地预订' },
  { id: '714', text: '网球教练', hotScore: 83, category: '运动健身', description: '网球培训' },
  { id: '715', text: '乒乓球培训', hotScore: 81, category: '运动健身', description: '球类培训' },
  
  // 汽车服务 (100条)
  { id: '801', text: '比亚迪汉EV', hotScore: 96, category: '汽车服务', description: '新能源汽车' },
  { id: '802', text: '特斯拉Model 3', hotScore: 94, category: '汽车服务', description: '电动汽车' },
  { id: '803', text: '小鹏P7', hotScore: 91, category: '汽车服务', description: '智能电车' },
  { id: '804', text: '蔚来ES6', hotScore: 89, category: '汽车服务', description: '纯电SUV' },
  { id: '805', text: '理想L9', hotScore: 87, category: '汽车服务', description: '增程式SUV' },
  
  { id: '806', text: '宝马5系', hotScore: 93, category: '汽车服务', description: '豪华轿车' },
  { id: '807', text: '奔驰E级', hotScore: 91, category: '汽车服务', description: '商务轿车' },
  { id: '808', text: '奥迪A6L', hotScore: 89, category: '汽车服务', description: '豪华中大型车' },
  { id: '809', text: '雷克萨斯ES', hotScore: 87, category: '汽车服务', description: '豪华轿车' },
  { id: '810', text: '凯迪拉克CT5', hotScore: 85, category: '汽车服务', description: '运动轿车' },
  
  { id: '811', text: '汽车保养', hotScore: 90, category: '汽车服务', description: '维修保养' },
  { id: '812', text: '轮胎更换', hotScore: 88, category: '汽车服务', description: '轮胎服务' },
  { id: '813', text: '汽车美容', hotScore: 86, category: '汽车服务', description: '美容服务' },
  { id: '814', text: '车险报价', hotScore: 84, category: '汽车服务', description: '保险服务' },
  { id: '815', text: '违章查询', hotScore: 82, category: '汽车服务', description: '车主服务' },
  
  // 金融理财 (50条)
  { id: '901', text: '余额宝收益', hotScore: 92, category: '金融理财', description: '货币基金' },
  { id: '902', text: '股票开户', hotScore: 90, category: '金融理财', description: '证券服务' },
  { id: '903', text: '基金定投', hotScore: 88, category: '金融理财', description: '投资理财' },
  { id: '904', text: '信用卡申请', hotScore: 86, category: '金融理财', description: '银行卡服务' },
  { id: '905', text: '贷款利率', hotScore: 84, category: '金融理财', description: '贷款服务' },
  { id: '906', text: '房贷计算器', hotScore: 87, category: '金融理财', description: '房贷工具' },
  { id: '907', text: '社保查询', hotScore: 85, category: '金融理财', description: '社保服务' },
  { id: '908', text: '公积金提取', hotScore: 83, category: '金融理财', description: '公积金服务' },
  { id: '909', text: '个人所得税', hotScore: 81, category: '金融理财', description: '税务服务' },
  { id: '910', text: '保险产品', hotScore: 79, category: '金融理财', description: '保险服务' },
  
  // 生活服务 (50条)
  { id: '1001', text: '外卖订餐', hotScore: 98, category: '生活服务', description: '餐饮外卖' },
  { id: '1002', text: '滴滴打车', hotScore: 96, category: '生活服务', description: '网约车' },
  { id: '1003', text: '快递查询', hotScore: 94, category: '生活服务', description: '快递服务' },
  { id: '1004', text: '家政服务', hotScore: 91, category: '生活服务', description: '家政清洁' },
  { id: '1005', text: '搬家公司', hotScore: 89, category: '生活服务', description: '搬家服务' },
  { id: '1006', text: '洗衣服务', hotScore: 86, category: '生活服务', description: '洗衣O2O' },
  { id: '1007', text: '开锁服务', hotScore: 84, category: '生活服务', description: '上门开锁' },
  { id: '1008', text: '水电维修', hotScore: 82, category: '生活服务', description: '维修服务' },
  { id: '1009', text: '宠物美容', hotScore: 88, category: '生活服务', description: '宠物服务' },
  { id: '1010', text: '代驾服务', hotScore: 85, category: '生活服务', description: '代驾' },
]

/**
 * 获取所有分类
 */
export function getCategories(): string[] {
  const categories = new Set<string>()
  mockSuggestions.forEach(item => {
    categories.add(item.category)
  })
  return Array.from(categories).sort()
}

/**
 * 获取分类统计
 */
export function getCategoryStats() {
  const stats = new Map<string, number>()
  mockSuggestions.forEach(item => {
    stats.set(item.category, (stats.get(item.category) || 0) + 1)
  })
  return Array.from(stats.entries()).map(([name, count]) => ({ name, count }))
}

