// 随机生成中文姓名
const lastNames = ['张', '高', '李'];
const firstNames = ['伟', '福', '丽', '秀英', '诺', '静', '朗', '强'];

export const generateRandomName = () => {
  const lastName = lastNames[Math.floor(Math.random() * lastNames.length)];
  const firstName = firstNames[Math.floor(Math.random() * firstNames.length)];
  return lastName + firstName;
};

// 随机生成地址
const cities = ['北京市', '海口市', '大连市', '深圳市'];
const districts = ['美兰区', '秀英区', '海甸岛'];
const streets = ['人民街道', '建国路', '复兴路' , '望京街'];
const communities = ['海南大学', '海南医科大学', '海南开放大学', '中科院'];

export const generateRandomAddress = () => {
  const city = cities[Math.floor(Math.random() * cities.length)];
  const district = districts[Math.floor(Math.random() * districts.length)];
  const street = streets[Math.floor(Math.random() * streets.length)];
  const community = communities[Math.floor(Math.random() * communities.length)];
  const building = Math.floor(Math.random() * 20 + 1);
  const unit = Math.floor(Math.random() * 6 + 1);
  const room = Math.floor(Math.random() * 2000 + 101);

  return `${city}${district}${street}${community}${building}号楼${unit}单元${room}室`;
};

// 随机生成面积（50-300平方米）
export const generateRandomArea = () => {
  return Number((Math.random() * (300 - 50) + 50).toFixed(2));
};

// 随机生成价格（50-1000万）
export const generateRandomPrice = () => {
  return Number((Math.random() * (10000000 - 500000) + 500000).toFixed(2));
}; 
