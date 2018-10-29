<template lang="pug">
  v-container
    v-layout(v-for="(actions, key) in groupedActions", :key="key", row)
      v-flex(xs5)
        v-layout(row)
          v-flex(xs5)
            div {{ key }}
          v-flex(xs7)
            div(v-for="action in actions", :key="`sidebar-${action._id}`") {{ action.desc }}
      v-flex(xs7)
        v-layout(row)
          v-flex(v-for="role in roles", :key="`role-${role._id}`")
            div {{ role._id }}
            div(v-for="action in actions", :key="`checkboxes-${action._id}`")
              v-checkbox(
              :input-value="getCheckboxValue(role, action._id)",
              @change="changeCheckboxValue($event, role, action)"
              )
    v-btn(@click="submit") Submit
</template>

<script>
import get from 'lodash/get';
import isEmpty from 'lodash/isEmpty';
import isUndefined from 'lodash/isUndefined';

import entitiesActionMixin from '@/mixins/entities/action';
import entitiesRoleMixin from '@/mixins/entities/role';

export default {
  mixins: [entitiesActionMixin, entitiesRoleMixin],
  data() {
    return {
      pending: true,
      changedRoles: {},
    };
  },
  computed: {
    groupedActions() {
      return this.actions.reduce((acc, action) => {
        if (action.id.startsWith('view') || action.id.startsWith('userview')) {
          acc.view.push(action);
        } else if (action.id.startsWith('models')) {
          acc.technical.push(action);
        } else {
          acc.business.push(action);
        }

        return acc;
      }, { business: [], view: [], technical: [] });
    },
    getCheckboxValue() {
      return (role, actionId) => {
        const changedValue = get(this.changedRoles, `${role._id}.${actionId}`);

        if (!isUndefined(changedValue)) {
          return changedValue;
        }

        if (role.rights[actionId]) {
          return Boolean(role.rights[actionId].checksum);
        }

        return false;
      };
    },
  },
  mounted() {
    this.fetchActionsList({
      params: { limit: 10 },
    });

    this.fetchRolesList({
      params: { limit: 10000 },
    });
  },
  methods: {
    submit() {
    },

    changeCheckboxValue(event, role, action) {
      if (!this.changedRoles[role._id]) {
        this.changedRoles[role._id] = {};
      }

      if (!isUndefined(this.changedRoles[role._id][action._id])) {
        this.$delete(this.changedRoles[role._id], action._id);

        if (isEmpty(this.changedRoles[role._id])) {
          this.$delete(this.changedRoles, role._id);
        }
      } else {
        this.$set(this.changedRoles[role._id], action._id, Number(event));
      }
    },
  },
};
</script>
