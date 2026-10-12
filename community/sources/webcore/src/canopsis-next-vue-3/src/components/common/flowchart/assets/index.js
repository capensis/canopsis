const modules = import.meta.glob('./**/*.svg', {
  eager: true,
  query: '?raw',
  import: 'default',
});

const assetGroups = Object.entries(modules).reduce((acc, [fileName, file]) => {
  const [, groupName] = fileName.match(/\.\/(.+)\/.+\.svg$/);

  if (acc[groupName]) {
    acc[groupName].push(file);
  } else {
    acc[groupName] = [file];
  }

  return acc;
}, {});

export default assetGroups;
