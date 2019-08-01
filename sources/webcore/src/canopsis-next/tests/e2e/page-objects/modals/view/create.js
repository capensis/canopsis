// https://nightwatchjs.org/guide/#working-with-page-objects

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');
const { DEFAULT_PAUSE } = require('../../../config');

const commands = {
  clearViewName() {
    return this.customClearValue('@viewFieldName');
  },
  clearViewTitle() {
    return this.customClearValue('@viewFieldTitle');
  },
  clearViewDescription() {
    return this.customClearValue('@viewFieldDescription');
  },
  clearViewGroupTags() {
    const { viewFieldGroupTagsChipsRemove } = this.elements;

    this.api.elements(
      viewFieldGroupTagsChipsRemove.locateStrategy,
      viewFieldGroupTagsChipsRemove.selector,
      ({ value = [] }) => {
        value.forEach(item => this.api.elementIdClick(item.ELEMENT).pause(DEFAULT_PAUSE));
      },
    );
    return this;
  },
  clearViewGroupIds() {
    return this.customClearValue('@viewFieldGroupIds');
  },
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
  clickViewDeleteButton() {
    return this.customClick('@viewDeleteButton');
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
      viewFieldGroupTagsChipsRemove: sel('.v-input.v-select--chips .v-select__selections .v-chip .v-chip__close'),
      viewFieldGroupIds: sel('viewFieldGroupIds'),
      viewSubmitButton: sel('viewSubmitButton'),
      viewDeleteButton: sel('viewDeleteButton'),
    }),
  },
  commands: [commands],
});
