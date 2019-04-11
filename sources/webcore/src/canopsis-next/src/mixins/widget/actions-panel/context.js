import { BUSINESS_USER_RIGHTS_ACTIONS_MAP } from '@/constants';

/**
 * @mixin Mixin for the alarms list actions panel, show modal of the action
 */
export default {
  methods: {
    actionsAccessFilterHandler({ type }) {
      const right = BUSINESS_USER_RIGHTS_ACTIONS_MAP.context[type];

      if (!right) {
        return true;
      }

      return this.checkAccess(right);
    },
  },
};
