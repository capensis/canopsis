import authMixin from '@/mixins/auth';
import entitiesRoleMixin from '@/mixins/entities/role';
import entitiesRightMixin from '@/mixins/entities/right';

export default {
  mixins: [authMixin, entitiesRoleMixin, entitiesRightMixin],
  methods: {
    async createRightAndAddIntoRole(right, checksum) {
      try {
        await this.createRight({ data: right });
        await this.addRightIntoRole(right._id, checksum);

        return this.fetchCurrentUser();
      } catch (err) {
        this.$popups.error({ text: this.$t('modals.view.errors.rightCreating') });

        return Promise.resolve();
      }
    },

    async removeRightAndRemoveRightFromRoleById(id) {
      try {
        await Promise.all([
          this.removeRight({ id }),
          this.removeRightFromRoleById(id),
        ]);

        return this.fetchCurrentUser();
      } catch (err) {
        this.$popups.error({ text: this.$t('modals.view.errors.rightRemoving') });

        return Promise.resolve();
      }
    },
  },
};
