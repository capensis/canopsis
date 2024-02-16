const fs = require('fs');
const path = require('path');

const imagesPath = path.resolve(process.cwd(), 'src', 'assets', 'images');

const images = fs.readdirSync(imagesPath);

module.exports = images.reduce((acc, fileName) => {
  const iconName = fileName.replace('.svg', '');

  acc[iconName] = {
    component: {
      name: fileName,
      template: fs.readFileSync(path.resolve(imagesPath, fileName)).toString(),
    },
  };

  return acc;
}, {});
