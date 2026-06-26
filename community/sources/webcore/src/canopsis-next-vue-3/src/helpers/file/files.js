export const saveFile = (data, filename) => {
  const blob = data instanceof Blob ? data : new Blob([data]);
  const link = document.createElement('a');

  link.href = URL.createObjectURL(blob);
  link.download = filename;
  link.click();

  URL.revokeObjectURL(link.href);
};
