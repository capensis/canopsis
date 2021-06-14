import { BUSINESS_USER_PERMISSIONS_ACTIONS_MAP } from '@/constants';

/**
 * @mixin Mixin for the alarms list actions panel, show modal of the action
 */
export default {
  methods: {
    actionsAccessFilterHandler({ type }) {
      const permission = BUSINESS_USER_PERMISSIONS_ACTIONS_MAP.context[type];

      if (!permission) {
        return true;
      }

      return this.checkAccess(permission);
    },
  },
};
