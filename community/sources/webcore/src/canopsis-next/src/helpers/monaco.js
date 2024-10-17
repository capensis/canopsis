import { isUndefined } from 'lodash';

/**
 * Determines the completion item kind for a given completion object.
 *
 * @param {Object} monaco - The Monaco editor instance, which provides access to language features.
 * @param {*} completion - The completion object whose type is to be determined. This can be of any type.
 * @param {boolean} isMember - A flag indicating whether the completion is a member of an object.
 * @returns {monaco.languages.CompletionItemKind} - Returns the appropriate `CompletionItemKind` based on the type of
 * the completion.
 *
 * - If the `completion` is an object, it returns `CompletionItemKind.Class`.
 * - If the `completion` is a function, it returns `CompletionItemKind.Method` if `isMember` is true, otherwise
 * `CompletionItemKind.Function`.
 * - For all other types, it returns `CompletionItemKind.Property` if `isMember` is true, otherwise
 * `CompletionItemKind.Variable`.
 */
export const getCompletionType = (monaco, completion, isMember) => {
  switch ((typeof completion).toLowerCase()) {
    case 'object':
      return monaco.languages.CompletionItemKind.Class;

    case 'function':
      return isMember
        ? monaco.languages.CompletionItemKind.Method
        : monaco.languages.CompletionItemKind.Function;

    default:
      return isMember
        ? monaco.languages.CompletionItemKind.Property
        : monaco.languages.CompletionItemKind.Variable;
  }
};

/**
 * Registers a JavaScript completion item provider for the Monaco editor.
 * This function enables code completion suggestions for JavaScript code within the editor.
 *
 * @param {Object} monaco - The Monaco editor instance. It is expected to have a `languages` property
 *                          with a `registerCompletionItemProvider` method.
 * @param {Object} completions - An object representing the available completions. The keys are the
 *                               names of the completions, and the values are the completion details.
 *                               The structure of this object should reflect the hierarchy of the
 *                               JavaScript objects and functions you want to provide completions for.
 *
 * @returns {Object|undefined} - Returns the result of the `registerCompletionItemProvider` call if
 *                               `monaco` is provided, otherwise returns `undefined`.
 *
 * The completion provider is triggered by specific characters (e.g., '.', '(') and provides
 * suggestions based on the current context in the editor. It analyzes the text before the cursor
 * to determine the appropriate completions to suggest.
 *
 * The completion items include:
 * - `label`: The text displayed in the suggestion list.
 * - `kind`: The type of completion item (e.g., Class, Method, Function, Property, Variable).
 * - `detail`: Additional information about the completion item, such as its type.
 * - `insertText`: The text to insert when the completion is selected.
 * - `documentation`: (Optional) Documentation or description of the completion item, extracted
 *                    from the function's string representation if applicable.
 *
 * The function handles member access by checking if the active typing ends with a dot ('.').
 * It navigates through the `completions` object hierarchy to find the relevant completions.
 * It also filters out properties that start with '__' and handles errors when accessing prototypes.
 */
export const registerJavaScriptCompletion = (monaco, completions) => monaco && (
  monaco.languages.registerCompletionItemProvider('javascript', {
    triggerCharacters: ['.'],

    provideCompletionItems: (model, position) => {
      const lastChars = model.getValueInRange({
        startLineNumber: position.lineNumber,
        startColumn: 0,
        endLineNumber: position.lineNumber,
        endColumn: position.column,
      });

      const words = lastChars.trim().split(' ');
      const activeTyping = words[words.length - 1];
      const path = activeTyping.split('.');
      const lastKey = path.pop();
      let current = completions;

      for (const key of path) {
        if (isUndefined(current[key])) {
          return undefined;
        }

        current = current[key];
      }

      const prefix = path.join('.');

      const suggestions = Object.entries(current).reduce((acc, [key, value]) => {
        if (!lastKey || key.startsWith(lastKey)) {
          const suggestion = {
            label: `${prefix}${prefix ? '.' : ''}${key}`,
            kind: getCompletionType(monaco, value, path.length > 0),
            // eslint-disable-next-line no-proto
            detail: value?.__proto__?.constructor?.name ?? typeof value,
            insertText: key,
          };

          if (suggestion.detail.toLowerCase() === 'function') {
            suggestion.insertText += '(';
            [suggestion.documentation] = value.toString().split('{');
          }

          acc.push(suggestion);
        }

        return acc;
      }, []);

      return {
        suggestions,
      };
    },
  })
);
