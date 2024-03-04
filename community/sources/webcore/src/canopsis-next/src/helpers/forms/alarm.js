/**
 * @typedef {Object} AlarmOldLink
 * @property {string} label
 * @property {string} link
 */

/**
 * @typedef {Object} AlarmLink
 * @property {string} icon_name
 * @property {string} label
 * @property {string} url
 * @property {boolean} [single]
 * @property {string} [rule_id]
 * @property {LinkRuleAction} [action]
 */

/**
 * @typedef {Object<string, AlarmLink[] | AlarmOldLink[]>} AlarmLinks
 */

/**
 * @typedef {Object} AlarmEvent
 * @property {number} val
 * @property {string} a
 * @property {number} t
 * @property {string} _t
 * @property {string} m
 */

/**
 * @typedef {Object} Correlation
 * @property {number} total
 * @property {Alarm} data
 */

/**
 * @typedef {Object} AlarmRule
 * @property {string} id
 * @property {string} name
 */

/**
 * @typedef {Object} AlarmAssignedInstructionExecution
 * @property {string} _id
 * @property {number} status
 */

/**
 * @typedef {Object} AlarmAssignedInstruction
 * @property {string} _id
 * @property {string} name
 * @property {?AlarmAssignedInstructionExecution} execution
 */

/**
 * @typedef {Object} AlarmValuePbehaviorInfo
 * @property {string} canonical_type
 * @property {string} icon_name
 * @property {string} id
 * @property {string} name
 * @property {string} reason
 * @property {string} type
 * @property {string} type_name
 * @property {number} timestamp
 */

/**
 * @typedef {Object} AlarmValue
 * @property {string[]} long_output_history
 * @property {number} last_event_date
 * @property {string} connector_name
 * @property {string} initial_long_output
 * @property {string} output
 * @property {string} initial_output
 * @property {string[]} children
 * @property {Object} extra
 * @property {number} last_update_date
 * @property {number} resolved
 * @property {string} resource
 * @property {number} creation_date
 * @property {string} display_name
 * @property {number} total_state_changes
 * @property {string[]} tags
 * @property {string} long_output
 * @property {string} component
 * @property {string} connector
 * @property {Object} infos_rule_version
 * @property {InfosObject} infos
 * @property {string[]} parents
 * @property {AlarmEvent} ack
 * @property {AlarmEvent} state
 * @property {AlarmEvent} status
 * @property {AlarmEvent} ticket
 * @property {AlarmEvent} snooze
 * @property {AlarmEvent} canceled
 * @property {AlarmEvent} last_comment
 * @property {AlarmEvent[]} [steps]
 * @property {AlarmValuePbehaviorInfo} [pbehavior_info]
 */

/**
 * @typedef {Pbehavior} AlarmPbehavior
 * @property {Comment} last_comment
 */

/**
 * @typedef {Object} Alarm
 * @property {string} _id
 * @property {Entity} entity
 * @property {boolean} is_meta_alarm
 * @property {AlarmAssignedInstruction[]} [assigned_instructions]
 * @property {number} [instruction_execution_icon]
 * @property {string[]} running_manual_instructions
 * @property {string[]} running_auto_instructions
 * @property {string[]} failed_manual_instructions
 * @property {string[]} failed_auto_instructions
 * @property {string[]} successful_manual_instructions
 * @property {string[]} successful_auto_instructions
 * @property {boolean} [children_instructions]
 * @property {InfosObject} infos
 * @property {AlarmRule} rule
 * @property {Correlation} consequences
 * @property {Correlation} causes
 * @property {AlarmPbehavior} pbehavior
 * @property {AlarmLinks} links
 * @property {number} t
 * @property {AlarmValue} v
 */
