package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nabilfikrisp/go-crud/internal/controller/restapi/v1/request"
	"github.com/nabilfikrisp/go-crud/internal/entity"
	"github.com/nabilfikrisp/go-crud/internal/usecase"
	"github.com/nabilfikrisp/go-crud/pkg/ptr"
)

func (r *V1) createContact(c *gin.Context) {
	var body request.CreateContact

	if err := c.ShouldBindJSON(&body); err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "restapi - v1 - createContact")
		validationErrorResponse(c, err)
		return
	}

	contact, err := r.c.Create(c.Request.Context(), usecase.CreateContact{
		FirstName:    body.FirstName,
		LastName:     body.LastName,
		Email:        body.Email,
		PhoneNumber:  body.PhoneNumber,
		Relationship: body.Relationship,
	})

	if err != nil {
		r.l.Error(err, "restapi - v1 - createContact")
		errorResponse(c, http.StatusInternalServerError, "Failed to create contact")
		return
	}

	c.JSON(http.StatusCreated, contact)
}

func (r *V1) getContactByID(c *gin.Context) {
	id := c.Param("id")

	contact, err := r.c.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, entity.ErrContactRelationshipInvalid) {
			errorResponse(c, http.StatusBadRequest, "Invalid relationship")
			return
		}

		if errors.Is(err, entity.ErrContactNotFound) {
			errorResponse(c, http.StatusNotFound, "Contact not found")
			return
		}

		r.l.Error(err, "restapi - v1 - getContactByID")
		errorResponse(c, http.StatusInternalServerError, "Failed to get contact")
		return
	}

	c.JSON(http.StatusOK, contact)
}

func (r *V1) listContacts(c *gin.Context) {
	filter := entity.ContactFilter{
		Limit:  ptr.Uint64(10),
		Offset: ptr.Uint64(0),
	}

	if v := c.Query("first_name"); v != "" {
		filter.FirstName = ptr.String(v)
	}

	if v := c.Query("last_name"); v != "" {
		filter.LastName = ptr.String(v)
	}

	if v := c.Query("email"); v != "" {
		filter.Email = ptr.String(v)
	}

	if v := c.Query("phone_number"); v != "" {
		filter.PhoneNumber = ptr.String(v)
	}

	if v := c.Query("relationship"); v != "" {
		rel := entity.ContactRelationship(v)
		if !rel.Valid() {
			errorResponse(c, http.StatusBadRequest, "invalid relationship")
			return
		}
		filter.Relationship = &rel
	}

	if v, err := strconv.ParseUint(c.Query("limit"), 10, 64); err == nil {
		filter.Limit = &v
	}
	if v, err := strconv.ParseUint(c.Query("offset"), 10, 64); err == nil {
		filter.Offset = &v
	}

	contacts, total, err := r.c.List(c.Request.Context(), filter)
	if err != nil {
		r.l.Error(err, "restapi - v1 - listContacts")
		errorResponse(c, http.StatusInternalServerError, "failed to list contacts")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"contacts": contacts,
		"total":    total,
	})
}

func (r *V1) updateContact(c *gin.Context) {
	id := c.Param("id")

	var body request.UpdateContact

	if err := c.ShouldBindJSON(&body); err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "restapi - v1 - updateContact")
		validationErrorResponse(c, err)
		return
	}

	contact, err := r.c.Update(c.Request.Context(), id, usecase.UpdateContact{
		FirstName:    body.FirstName,
		LastName:     body.LastName,
		Email:        body.Email,
		PhoneNumber:  body.PhoneNumber,
		Relationship: body.Relationship,
	})
	if err != nil {
		if errors.Is(err, entity.ErrContactNotFound) {
			errorResponse(c, http.StatusNotFound, "Contact not found")
			return
		}

		if errors.Is(err, entity.ErrContactRelationshipInvalid) {
			errorResponse(c, http.StatusBadRequest, "Invalid relationship")
			return
		}

		r.l.Error(err, "restapi - v1 - updateContact")
		errorResponse(c, http.StatusInternalServerError, "Failed to update contact")
		return
	}

	c.JSON(http.StatusOK, contact)
}

func (r *V1) deleteContact(c *gin.Context) {
	id := c.Param("id")

	if err := r.c.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, entity.ErrContactNotFound) {
			errorResponse(c, http.StatusNotFound, "Contact not found")
			return
		}

		r.l.Error(err, "restapi - v1 - deleteContact")
		errorResponse(c, http.StatusInternalServerError, "Failed to delete contact")
		return
	}

	c.Status(http.StatusNoContent)
}
