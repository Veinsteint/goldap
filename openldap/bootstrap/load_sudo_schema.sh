#!/bin/bash
# Load sudo schema into OpenLDAP container
# This script runs from the host and executes commands inside the container

set -e

CONTAINER_NAME="${CONTAINER_NAME:-goldap.ldap}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LDIF_FILE="${SCRIPT_DIR}/ldif/sudo-schema.ldif"

echo "=========================================="
echo "LDAP Sudo Schema Loader"
echo "=========================================="
echo "Container: ${CONTAINER_NAME}"
echo ""

# Check if container exists
echo "1. Checking container..."
if ! docker ps --format "{{.Names}}" | grep -q "^${CONTAINER_NAME}$"; then
    echo "[ERROR] Container ${CONTAINER_NAME} not found or not running"
    exit 1
fi
echo "[OK] Container is running"

# Copy LDIF file to container
echo "2. Copying schema file..."
docker cp "${LDIF_FILE}" ${CONTAINER_NAME}:/tmp/sudo-schema.ldif
echo "[OK] File copied"

# Check if schema already exists
echo "3. Checking existing schema..."
if docker exec ${CONTAINER_NAME} ldapsearch -Y EXTERNAL -H ldapi:/// -b "cn=schema,cn=config" "(cn=*sudo*)" dn 2>&1 | grep -q "sudo"; then
    echo "[SKIP] Sudo schema already loaded"
    exit 0
fi
echo "[OK] Schema not found, will load"

# Load schema
echo "4. Loading sudo schema..."
docker exec ${CONTAINER_NAME} ldapadd -Y EXTERNAL -H ldapi:/// -f /tmp/sudo-schema.ldif 2>&1 || {
    EXIT_CODE=$?
    if [ ${EXIT_CODE} -eq 68 ]; then
        echo "[OK] Schema already exists"
    else
        echo "[ERROR] Failed to load schema (exit code: ${EXIT_CODE})"
        exit 1
    fi
}
echo "[OK] Schema loaded"

# Verify
echo "5. Verifying..."
if docker exec ${CONTAINER_NAME} ldapsearch -Y EXTERNAL -H ldapi:/// -b "cn=schema,cn=config" "(cn=*sudo*)" dn 2>&1 | grep -q "sudo"; then
    echo "[OK] Sudo schema verified"
else
    echo "[ERROR] Verification failed"
    exit 1
fi

echo ""
echo "=========================================="
echo "[OK] Sudo schema ready"
echo "=========================================="

