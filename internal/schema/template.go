package schema

// DefaultTemplate is the template for a new schema.yaml file
const DefaultTemplate = `# pgmigrate schema definition
# Documentation: https://github.com/matroidbe/pgmigrate

# Schemas managed by pgmigrate - only tables in these schemas will be tracked
managed_schemas:
  - public

# Tables to create/manage
tables:
  public.example:
    columns:
      - name: id
        type: bigserial
        primary_key: true
      - name: name
        type: text
        not_null: true
      - name: email
        type: varchar(255)
      - name: active
        type: boolean
        default: "true"
      - name: created_at
        type: timestamptz
        not_null: true
        default: "now()"
      - name: updated_at
        type: timestamptz

    indexes:
      - name: example_email_idx
        columns: [email]
        unique: true
`
