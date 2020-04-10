/**
 * Function for create custom variables button.
 *
 * @param {Array} data
 * @param {String} containerClassName
 * @param {String} icon
 * @return {Object}
 */
export const createJoditVariablesButton = (data = [], { containerClassName, icon = 'source' }) => ({
  name: 'variables',
  icon,
  data,
  popup(editor, _, control) {
    const popup = document.createElement('div');
    popup.classList.add(containerClassName);

    control.data.forEach(({ label, value }) => {
      const button = document.createElement('button');
      button.type = 'button';
      button.innerHTML = label;

      button.addEventListener('click', () => editor.selection.insertHTML(value));

      popup.appendChild(button);
    });

    return popup;
  },
});
