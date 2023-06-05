export const triggerDocumentEvent = (event) => {
  document.dispatchEvent(event);
};

export const triggerWindowEvent = (event) => {
  window.dispatchEvent(event);
};

export const triggerDocumentMouseEvent = (type, data) => {
  triggerDocumentEvent(new MouseEvent(type, data));
};

export const triggerDocumentKeyboardEvent = (type, data) => {
  triggerDocumentEvent(new KeyboardEvent(type, data));
};

export const triggerWindowScrollEvent = (detail) => {
  triggerWindowEvent(new CustomEvent('scroll', { detail }));
};

export const triggerWindowKeyboardEvent = (type, data) => {
  triggerWindowEvent(new KeyboardEvent(type, data));
};
