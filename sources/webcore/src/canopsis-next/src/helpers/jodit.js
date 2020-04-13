/**
 * Function for create custom variables button.
 *
 * @param {Function} onClick
 * @return {Object}
 */
export const createJoditVariablesButton = ({ onClick }) => ({
  name: 'variables',
  mode: 3,
  getContent: (editor) => {
    const controlButton = document.createElement('span');

    controlButton.classList.add('jodit-variables-help-btn');
    controlButton.addEventListener('click', event => onClick(editor, event));

    return controlButton;
  },
});
