// https://nightwatchjs.org/guide/#working-with-page-objects

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const commands = {
  setViewName(value) {
    return this.customSetValue('@viewFieldName', value);
  },
  setViewTitle(value) {
    return this.customSetValue('@viewFieldTitle', value);
  },
  setViewDescription(value) {
    return this.customSetValue('@viewFieldDescription', value);
  },
  setViewGroupTags(value) {
    return this.customSetValue('@viewFieldGroupTags', value)
      .customKeyup('@viewFieldGroupTags', 'ENTER');
  },
  setViewGroupIds(value) {
    return this.customSetValue('@viewFieldGroupIds', value)
      .customKeyup('@viewFieldGroupIds', 'ENTER');
  },
  clickViewEnabled() {
    return this.customClick('@viewFieldEnabled');
  },
  clickViewSubmitButton() {
    return this.customClick('@viewSubmitButton');
  },
};

const modalSelector = sel('createViewModal');

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      viewFieldName: sel('viewFieldName'),
      viewFieldTitle: sel('viewFieldTitle'),
      viewFieldDescription: sel('viewFieldDescription'),
      viewFieldEnabled: `.v-input${sel('viewFieldEnabled')} .v-input--selection-controls__ripple`,
      viewFieldGroupTags: sel('viewFieldGroupTags'),
      viewFieldGroupIds: sel('viewFieldGroupIds'),
      viewSubmitButton: sel('viewSubmitButton'),
    }),
  },
  commands: [commands],
});
