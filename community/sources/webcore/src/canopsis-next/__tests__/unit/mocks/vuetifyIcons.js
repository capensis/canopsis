import fs from 'node:fs';
import path from 'node:path';

const imagesPath = path.resolve(process.cwd(), 'src', 'assets', 'images');

const images = fs.readdirSync(imagesPath);

export default images.reduce((acc, fileName) => {
  const iconName = fileName.replace('.svg', '');

  acc[iconName] = {
    component: {
      name: fileName,
      template: fs.readFileSync(path.resolve(imagesPath, fileName)).toString(),
    },
  };

  return acc;
}, {});
