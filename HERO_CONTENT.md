# Hero Content Feature

The SCORM Quiz Generator now supports **hero content** - contextual content displayed prominently above questions. This is perfect for:

- 🎣 **Phishing Training** - Show realistic email examples
- 💻 **Code Review** - Display code snippets
- 📄 **Document Analysis** - Present text passages
- 🖼️ **Visual Context** - Include images

## How It Works

### For Users

#### When Generating Quizzes
When you generate a quiz on topics like "Phishing Awareness" or "Code Security", the AI will automatically include relevant hero content where appropriate.

#### Manual Editing
For any question, click the **✏️ Edit** button to:
1. Select hero content type from dropdown:
   - **None** - No hero content
   - **Text/Code** - Plain text or code snippets
   - **Email** - Email messages (styled with email formatting)
   - **HTML** - Custom HTML content
   - **Image URL** - Image from a URL
2. Enter the content in the text area
3. Save the question

### Hero Content Types

#### 1. Email (`heroType: "email"`)
Perfect for phishing training:

```json
{
  "heroContent": "From: IT Support <support@company-security.com>\nTo: you@company.com\nSubject: URGENT: Verify Your Account\n\nDear Employee,\n\nYour account will be suspended in 24 hours unless you verify your credentials immediately.\n\nClick here to verify: http://verify-account.suspicious-domain.com\n\nRegards,\nIT Security Team",
  "heroType": "email"
}
```

**Displays as:** Formatted email with gray background and purple accent border

#### 2. Text/Code (`heroType: "text"`)
For code snippets or text passages:

```json
{
  "heroContent": "function processPayment(amount) {\n  // TODO: Add validation\n  executeTransaction(amount);\n  return true;\n}",
  "heroType": "text"
}
```

**Displays as:** Monospaced code block with syntax-friendly formatting

#### 3. HTML (`heroType: "html"`)
For custom formatted content:

```json
{
  "heroContent": "<div style='color: red;'><strong>Warning:</strong> This is a test</div>",
  "heroType": "html"
}
```

**Displays as:** Rendered HTML (use with caution)

#### 4. Image (`heroType: "image"`)
For visual context:

```json
{
  "heroContent": "https://example.com/suspicious-email-screenshot.png",
  "heroType": "image"
}
```

**Displays as:** Responsive image

## Example Use Cases

### Phishing Awareness Quiz

**Subject:** "Email Phishing Detection"

**Question with Hero Content:**
```json
{
  "id": 1,
  "heroContent": "From: paypal@secure-payment-verify.com\nSubject: Your account has been limited\n\nDear PayPal User,\n\nWe detected unusual activity on your account. Click below to restore access:\nhttp://paypal-verify.tk/restore\n\nPayPal Security",
  "heroType": "email",
  "question": "What is the MOST suspicious indicator in this email?",
  "options": [
    "The email mentions unusual activity",
    "The sender domain is 'secure-payment-verify.com' not 'paypal.com'",
    "The email is too short",
    "It mentions account restoration"
  ],
  "answer": 1
}
```

### Code Security Quiz

**Subject:** "Secure Coding Practices"

**Question with Hero Content:**
```json
{
  "id": 1,
  "heroContent": "app.get('/user/:id', (req, res) => {\n  const query = `SELECT * FROM users WHERE id = ${req.params.id}`;\n  db.query(query, (err, results) => {\n    res.json(results);\n  });\n});",
  "heroType": "text",
  "question": "What security vulnerability exists in this code?",
  "options": [
    "Missing error handling",
    "SQL injection vulnerability",
    "Incorrect HTTP method",
    "No authentication"
  ],
  "answer": 1
}
```

## API Changes

### Question Model
```go
type Question struct {
    ID          int      `json:"id"`
    Question    string   `json:"question"`
    Options     []string `json:"options"`
    Answer      int      `json:"answer"`
    HeroContent string   `json:"heroContent,omitempty"` // NEW
    HeroType    string   `json:"heroType,omitempty"`    // NEW
}
```

### Creating Questions with Hero Content

```javascript
// Via Web UI API
fetch('/api/update-question', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    sessionId: 'session_123',
    questionIndex: 0,
    question: {
      id: 1,
      question: "Is this email legitimate?",
      options: ["Yes", "No", "Unsure", "Need more info"],
      answer: 1,
      heroContent: "From: CEO@company.co...",
      heroType: "email"
    }
  })
});
```

## SCORM Output

When you generate a SCORM package with hero content:

1. **Hero content displays first** - Above the question text
2. **Styled appropriately** - Based on heroType
3. **Fully responsive** - Works on all devices
4. **SCORM compliant** - No impact on LMS compatibility

## Best Practices

### ✅ Do:
- Use hero content for context-dependent questions
- Keep email examples realistic but safe (don't use real domains)
- Format code with proper indentation
- Keep content concise and relevant

### ❌ Don't:
- Include real credentials or sensitive data
- Use hero content for every question (only when needed)
- Embed malicious content or actual phishing attempts
- Use overly long content (breaks mobile experience)

## Tips for Phishing Training

### Realistic but Safe Emails

**Good Example:**
```
From: security@company-verify-system.com
Subject: Account Verification Required

Your account will be locked in 2 hours. Click to verify:
http://verify-portal.suspicious-tld.com
```

**Things to Include:**
- Misspelled domains (company-security vs company.security)
- Suspicious TLDs (.tk, .ml, etc.)
- Urgency tactics ("in 2 hours", "immediately")
- Generic greetings ("Dear User" instead of name)
- Grammar/spelling mistakes
- Mismatched sender addresses

### Question Patterns

1. **Identification:** "What makes this email suspicious?"
2. **Priority:** "What is the MOST concerning indicator?"
3. **Action:** "What should you do if you receive this?"
4. **Technical:** "Which technical indicator reveals this is phishing?"

## Browser Compatibility

Hero content works in:
- ✅ Chrome/Edge (Chromium)
- ✅ Firefox
- ✅ Safari
- ✅ Mobile browsers

## Migration

Existing quizzes without hero content:
- ✅ Continue to work perfectly
- ✅ Can be edited to add hero content
- ✅ Optional fields, fully backward compatible

## Support

For issues or questions about hero content:
1. Check that `heroType` matches content format
2. Verify JSON is valid
3. Test in SCORM preview before deploying
4. Check browser console for errors

---

**Happy quiz creating! 🎓**
