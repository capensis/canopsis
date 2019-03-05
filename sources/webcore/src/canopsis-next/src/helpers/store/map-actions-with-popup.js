import mapActionsWith from './map-actions-with';

/**
 * Mapping actions with popups
 *
 * @param actions
 */
export default function (actions) {
  return mapActionsWith(actions, async function success(popups) {
    let successPopup;

    if (popups) {
      successPopup = popups.success ? popups.success : popups;
    } else {
      successPopup = { text: this.$t('success.default') };
    }

    if (this.addSuccessPopup) {
      await this.addSuccessPopup(successPopup);
    }
  }, async function error(popups) {
    const errorPopup = popups && popups.error ? popups.error : { text: this.$t('errors.default') };

    if (this.addErrorPopup) {
      await this.addErrorPopup(errorPopup);
    }
  });
}
