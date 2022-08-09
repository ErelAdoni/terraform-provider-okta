package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

var (
	translationResource = &schema.Resource{
		Schema: map[string]*schema.Schema{
			"language": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subject": {
				Type:     schema.TypeString,
				Required: true,
			},
			"template": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
	// NOTE: update values from GET {{url}}/api/v1/templates/emails
	validEmailTemplateTypes = []string{
		"email.accountLockout",
		"email.ad.forgotPassword",
		"email.ad.forgotPasswordReset",
		"email.ad.selfServiceUnlock",
		"email.ad.welcome",
		"email.adminUsersListScheduledLifecycleStatusChange",
		"email.authenticatorEnrollment",
		"email.authenticatorReset",
		"email.automation",
		"email.emailActivation",
		"email.emailChangeConfirmation",
		"email.emailLinkAuthenticationTransaction",
		"email.emailLinkFactorVerification",
		"email.emailLinkRecoveryAdResetFactor",
		"email.emailLinkRecoveryAdResetPwdFactor",
		"email.emailLinkRecoveryAdUnlockFactor",
		"email.emailLinkRecoveryAdUnlockPwdFactor",
		"email.emailLinkRecoveryLdapResetPwdFactor",
		"email.emailLinkRecoveryLdapUnlockPwdFactor",
		"email.emailLinkRecoveryResetFactor",
		"email.emailLinkRecoveryResetPwdFactor",
		"email.emailLinkRecoveryUnlockFactor",
		"email.emailLinkRecoveryUnlockPwdFactor",
		"email.emailNewAlreadyChangedNotification",
		"email.emailNewChangeNotification",
		"email.emailTransactionVerification",
		"email.endUserScheduledLifecycleStatusChange",
		"email.factorEnrollment",
		"email.factorReset",
		"email.forgotPassword",
		"email.forgotPasswordDenied",
		"email.idpMyAccountChangeConfirmation",
		"email.passwordChanged",
		"email.pushVerifyActivation",
		"email.registrationActivation",
		"email.registrationEmailVerification",
		"email.selfServiceUnlock",
		"email.selfServiceUnlockOnUnlockedAccount",
		"email.signInFromNewDevice",
		"email.sunone.forgotPassword",
		"email.sunone.forgotPasswordDenied",
		"email.sunone.selfServiceUnlock",
		"email.sunone.welcome",
		"email.tempPassword",
		"email.welcome",
	}
)

func resourceTemplateEmail() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "Resource okta_template_email utilizes a private Okta API whose behavior may change or even removed. Resource okta_template_emal has been replaced by resource okta_email_customization which is supported by public Okta API.",
		CreateContext:      resourceTemplateEmailCreate,
		ReadContext:        resourceTemplateEmailRead,
		UpdateContext:      resourceTemplateEmailUpdate,
		DeleteContext:      resourceTemplateEmailDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"default_language": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "en",
			},
			"type": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Email template type",
				ForceNew:         true,
				ValidateDiagFunc: elemInSlice(validEmailTemplateTypes),
			},
			"translations": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     translationResource,
			},
		},
	}
}

func resourceTemplateEmailCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	temp := buildEmailTemplate(d)
	id := d.Get("type").(string)
	_, _, err := getAPISupplementFromMetadata(m).CreateEmailTemplate(ctx, *temp, nil)
	if err != nil {
		return diag.Errorf("failed to create email template: %v", err)
	}
	d.SetId(id)
	return resourceTemplateEmailRead(ctx, d, m)
}

func resourceTemplateEmailRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	temp, resp, err := getAPISupplementFromMetadata(m).GetEmailTemplate(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get email template: %v", err)
	}
	if temp == nil {
		d.SetId("")
		return nil
	}
	if temp.Id == "default" {
		d.SetId("")
		return nil
	}
	_ = d.Set("translations", flattenEmailTranslations(temp.Translations))
	_ = d.Set("default_language", temp.DefaultLanguage)
	return nil
}

func resourceTemplateEmailUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	temp := buildEmailTemplate(d)
	_, _, err := getAPISupplementFromMetadata(m).UpdateEmailTemplate(ctx, d.Id(), *temp, nil)
	if err != nil {
		return diag.Errorf("failed to update email template: %v", err)
	}
	return resourceTemplateEmailRead(ctx, d, m)
}

func resourceTemplateEmailDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_, err := getAPISupplementFromMetadata(m).DeleteEmailTemplate(ctx, d.Id())
	if err != nil {
		return diag.Errorf("failed to delete email template: %v", err)
	}
	return nil
}

func buildEmailTemplate(d *schema.ResourceData) *sdk.EmailTemplate {
	trans := map[string]*sdk.EmailTranslation{}
	rawTransList := d.Get("translations").(*schema.Set)

	for _, val := range rawTransList.List() {
		rawTrans := val.(map[string]interface{})
		trans[rawTrans["language"].(string)] = &sdk.EmailTranslation{
			Subject:  rawTrans["subject"].(string),
			Template: rawTrans["template"].(string),
		}
	}
	defaultLang := d.Get("default_language").(string)

	return &sdk.EmailTemplate{
		DefaultLanguage: defaultLang,
		Name:            "Custom",
		Type:            d.Get("type").(string),
		Translations:    trans,
		Subject:         trans[defaultLang].Subject,
		Template:        trans[defaultLang].Template,
	}
}

func flattenEmailTranslations(temp map[string]*sdk.EmailTranslation) *schema.Set {
	var rawSet []interface{}
	for key, val := range temp {
		rawSet = append(rawSet, map[string]interface{}{
			"language": key,
			"subject":  val.Subject,
			"template": val.Template,
		})
	}
	return schema.NewSet(schema.HashResource(translationResource), rawSet)
}
