package alarmaction

import "github.com/go-playground/validator/v10"

func ValidateCommentRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(CommentRequest)

	if len(r.StructuredComment) > 0 && r.Comment != "" {
		sl.ReportError(r.StructuredComment, "StructuredComment", "StructuredComment", "required_not_both", "Comment")
	}

	if len(r.StructuredComment) == 0 && r.Comment == "" {
		sl.ReportError(r.StructuredComment, "StructuredComment", "StructuredComment", "required_or", "Comment")
		sl.ReportError(r.Comment, "Comment", "Comment", "required_or", "StructuredComment")
	}
}
