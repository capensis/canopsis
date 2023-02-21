function isStringStart(chr) {
  const SINGLE_QUOTE_CODE = 0x22;
  const DOUBLE_QUOTE_CODE = 0x27;

  return chr === SINGLE_QUOTE_CODE || chr === DOUBLE_QUOTE_CODE;
}

/**
 * This function is parsing v-field original value to real path
 * Example: form.field[key] => ['field', key]
 *
 * @param {string} value - original path from v-field (v-field="form.field[key]")
 * @returns {Array}
 */
function parseField(value) {
  const START_BRACKET_CODE = 0x5B;
  const END_BRACKET_CODE = 0x5D;
  const POINT_CODE = 0x2e;
  const path = [];
  let hasString = false;
  let isPoint = true;
  let startPos = 0;
  let endPos;

  for (let index = 0; index < value.length; index += 1) {
    const chr = value.charCodeAt(index);

    /**
     * If there is '[' char symbol. It is saying us that it's start of field
     */
    if (chr === START_BRACKET_CODE) {
      /**
       * If we already have path. It will happened if we had point field.
       * Example: `form.field[key]` this condition will happened on '[' char symbol and we will save 'field'.
       * And it will not happened for `form['field'][key]`
       */
      if (startPos < index - 1) {
        path.push(JSON.stringify(value.slice(startPos, index)));
      }

      hasString = false;
      isPoint = false;
      startPos = index + 1;

      /**
       * We need to check it because if there is string path we will wrap value to JSON.stringify
       * Example: form['field'] or form["field"]
       */
      if (isStringStart(value.charCodeAt(startPos))) {
        hasString = true;
      }
      /**
       * If there is ']' char symbol. It is saying us that it's end of brackets field
       */
    } else if (chr === END_BRACKET_CODE) {
      endPos = index - 1;

      /**
       * If we have string inside brackets we will wrap it into JSON.stringify
       */
      if (hasString) {
        path.push(JSON.stringify(value.slice(startPos + 1, endPos)));
        hasString = false;
      } else {
        path.push(value.slice(startPos, endPos + 1));
      }

      startPos = index + 1;
      /**
       * If there is '.' char symbol. It is saying us that it's start of point field
       */
    } else if (chr === POINT_CODE && !hasString) {
      /**
       * If we had point in previous field we will save it
       * Example:
       * `form.field.anotherField` this condition will happened on second '.' char symbol and we will save 'field'
       */
      if (isPoint) {
        path.push(JSON.stringify(value.slice(startPos, index)));

        startPos = index + 1;
      }

      startPos = index + 1;
      isPoint = true;
    }
  }

  /**
   * If we had point field at the last position we will save it here
   * Example: `form.field.anotherField` -> 'anotherField'
   */
  if (isPoint) {
    path.push(JSON.stringify(value.slice(startPos)));
  }

  return path;
}

function updateFieldDirective(el, dir) {
  const { value, modifiers = {} } = dir;
  const {
    number,
    trim,
    mutate,
    model,
  } = modifiers;

  const path = parseField(value.trim());

  path.shift();

  const baseValueExpression = '$$v';
  const basePreviousPathExpression = '$$p';
  let valueExpression = baseValueExpression;

  if (trim) {
    valueExpression = `(typeof ${baseValueExpression} === 'string'`
      + `? ${baseValueExpression}.trim()`
      + `: ${baseValueExpression})`;
  }

  if (number) {
    valueExpression = `_n(${valueExpression})`;
  }

  const pathExpression = mutate ? `[${path}].concat(${basePreviousPathExpression} || [])` : `[${path}]`;
  const assignment = model
    ? `$updateFieldModel(${value}, ${basePreviousPathExpression}, ${valueExpression})`
    : `$updateField(${pathExpression}, ${valueExpression}, ${mutate})`;

  // eslint-disable-next-line no-param-reassign
  el.model = {
    value: `(${value})`,
    expression: JSON.stringify(value),
    callback: `function (${baseValueExpression}, ${basePreviousPathExpression}) {${assignment}}`,
  };
}

module.exports = updateFieldDirective;
