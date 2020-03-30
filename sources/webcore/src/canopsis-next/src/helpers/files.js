import FileSaver from 'file-saver';


export const saveFile = (blob, name) => FileSaver.saveAs(blob, name);
