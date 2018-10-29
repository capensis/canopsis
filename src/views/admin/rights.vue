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
      return (role, actionId, rightType = 15) => {
        const checkSum = get(role, `rights.${actionId}.checksum`, 0);
        const changedCheckSum = get(this.changedRoles, `${role._id}.${actionId}`);

        const currentCheckSum = isUndefined(changedCheckSum) ? checkSum : changedCheckSum;
        const actionRightType = currentCheckSum & rightType;

        return actionRightType === rightType;
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
      // console.log(this.changedRoles);
    },

    changeCheckboxValue(value, role, action, rightType = 15) {
      const factor = value ? 1 : -1;

      if (!this.changedRoles[role._id]) {
        this.changedRoles[role._id] = {};
      }

      if (!isUndefined(this.changedRoles[role._id][action._id])) {
        const nextCheckSum = this.changedRoles[role._id][action._id] + (factor * rightType);

        if (!nextCheckSum) {
          this.$delete(this.changedRoles[role._id], action._id);
        } else {
          this.$set(this.changedRoles[role._id], action._id, nextCheckSum);
        }

        if (isEmpty(this.changedRoles[role._id])) {
          this.$delete(this.changedRoles, role._id);
        }
      } else {
        const currentCheckSum = get(role.rights, `${action._id}.checksum`, 0);
        const nextCheckSum = currentCheckSum + (factor * rightType);

        this.$set(this.changedRoles[role._id], action._id, nextCheckSum);
      }
    },
  },
};
</script>
