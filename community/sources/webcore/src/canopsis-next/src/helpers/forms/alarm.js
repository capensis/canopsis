/**
 * @typedef {Object} AlarmLink
 * @property {string} label
 * @property {string} link
 */

/**
 * @typedef {Object.<string, AlarmLink[] | string[]>} AlarmLinks
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
 * @property {string} status
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
 * @property {string[]} filtered_children
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
 * @property {AlarmEvent} lastComment
 * @property {AlarmEvent[]} [steps]
 * @property {AlarmValuePbehaviorInfo} [pbehavior_info]
 */

/**
 * @typedef {Object} Alarm
 * @property {string} _id
 * @property {Entity} entity
 * @property {boolean} metaalarm
 * @property {AlarmAssignedInstruction[]} [assigned_instructions]
 * @property {boolean} [is_auto_instruction_running]
 * @property {boolean} [is_manual_instruction_waiting_result]
 * @property {boolean} [is_all_auto_instructions_completed]
 * @property {boolean} [children_instructions]
 * @property {InfosObject} infos
 * @property {AlarmRule} rule
 * @property {Correlation} consequences
 * @property {Correlation} causes
 * @property {Pbehavior} pbehavior
 * @property {string} links
 * @property {number} t
 * @property {AlarmValue} v
 */
