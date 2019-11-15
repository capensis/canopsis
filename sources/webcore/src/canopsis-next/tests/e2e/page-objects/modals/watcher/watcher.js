// https://nightwatchjs.org/guide/#working-with-page-objects

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const modalSelector = sel('watcherModal');

const commands = {
  clickRefreshButton() {
    return this.customClick('@refreshButton');
  },

  clickHideButton() {
    return this.customClick('@hideButton');
  },

  clickWatcherEntity(watcherId) {
    return this.customClick(this.el('@watcherEntity', watcherId));
  },

  clickWatcherEntityEditPbehaviors(watcherId) {
    return this.customClick(this.el('@watcherEntityEditPbehaviors', watcherId));
  },

  clickWatcherEntityAction(watcherId, eventType) {
    return this.customClick(this.el('@watcherEntityAction', watcherId, eventType));
  },
};

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      submitButton: sel('watcherSubmitButton'),
      cancelButton: sel('watcherCancelButton'),
    }),
    refreshButton: sel('watcherRefreshButton'),
    hideButton: sel('watcherHideButton'),
    watcherEntity: sel('entity-%s'),
    watcherEntityAction: `${sel('entity-%s')} ${sel('entityAction-%s')}`,
    watcherEntityEditPbehaviors: `${sel('entity-%s')} ${sel('entityActionEditPbehaviors')}`,
  },
  commands: [commands],
});
