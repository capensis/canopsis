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
              div(v-if="key === 'technical'")
                v-tooltip
                  v-checkbox(
                  slot="activator",
                  :input-value="getCheckboxValue(role, action._id, $constants.ACTIONS_RIGHTS_TYPES.create)",
                  @change="changeCheckboxValue($event, role, action, $constants.ACTIONS_RIGHTS_TYPES.create)"
                  )
                  span Create
                v-tooltip
                  v-checkbox(
                  slot="activator",
                  :input-value="getCheckboxValue(role, action._id, $constants.ACTIONS_RIGHTS_TYPES.read)",
                  @change="changeCheckboxValue($event, role, action, $constants.ACTIONS_RIGHTS_TYPES.read)"
                  )
                  span Read
                v-tooltip
                  v-checkbox(
                  slot="activator",
                  :input-value="getCheckboxValue(role, action._id, $constants.ACTIONS_RIGHTS_TYPES.update)",
                  @change="changeCheckboxValue($event, role, action, $constants.ACTIONS_RIGHTS_TYPES.update)"
                  )
                  span Update
                v-tooltip
                  v-checkbox(
                  slot="activator",
                  :input-value="getCheckboxValue(role, action._id, $constants.ACTIONS_RIGHTS_TYPES.delete)",
                  @change="changeCheckboxValue($event, role, action, $constants.ACTIONS_RIGHTS_TYPES.delete)"
                  )
                  span Delete
              div(v-else)
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
      return (role, actionId, rightType = this.$constants.ACTIONS_RIGHTS_TYPES.all) => {
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
      params: { limit: 70 },
    });

    this.fetchRolesList({
      params: { limit: 10000 },
    });
  },
  methods: {
    submit() {
      // console.log(this.changedRoles);
    },

    changeCheckboxValue(value, role, action, rightType = this.$constants.ACTIONS_RIGHTS_TYPES.all) {
      const currentCheckSum = get(role.rights, `${action._id}.checksum`, 0);
      const factor = value ? 1 : -1;

      if (!this.changedRoles[role._id]) {
        this.changedRoles[role._id] = {};
      }

      if (!isUndefined(this.changedRoles[role._id][action._id])) {
        const nextCheckSum = this.changedRoles[role._id][action._id] + (factor * rightType);

        if (currentCheckSum === nextCheckSum) {
          this.$delete(this.changedRoles[role._id], action._id);
        } else {
          this.$set(this.changedRoles[role._id], action._id, nextCheckSum);
        }

        if (isEmpty(this.changedRoles[role._id])) {
          this.$delete(this.changedRoles, role._id);
        }
      } else {
        const nextCheckSum = currentCheckSum + (factor * rightType);

        this.$set(this.changedRoles[role._id], action._id, nextCheckSum);
      }
    },
  },
};
</script>
