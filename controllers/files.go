package controllers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/httperror"
)

type Files struct {
	repo models.RepositoryService
	file models.FileService
}

func NewFiles(r models.RepositoryService, f models.FileService) *Files {
	return &Files{
		repo: r,
		file: f,
	}
}

func (h *Files) Download(c *Ctx) error {
	fileID := c.Path("fileID")
	if err := h.repo.AddFileDownload(c.Context(), fileID); err != nil {
		return err
	}
	return h.file.Get(c.Context(), fileID, c.Res)
}

func (h *Files) ConfirmDelete(c *Ctx) error {
	fileID := c.Path("fileID")

	//TODO: what if someone deleted the file in another tab?
	file, err := h.repo.FileByID(c.Context(), fileID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "file error: %s\n", err)
		/*
			TODO: strange issue: if the template exists, and this line
			is used, then this template is always used, even for the line
			at the bottom of the request????
		*/
		/*return c.HTML(
		http.StatusOK,
		"modal",
		"modals/confirm_file_was_deleted", Map{})*/
		return err
	}

	if !c.IsSpaceAdmin(c.User, file.Folder.Space) {
		return httperror.Forbidden
	}

	return c.HTML(
		http.StatusOK,
		"modal",
		"modals/confirm_delete_file", Map{
			"file": file,
		})
}

func (h *Files) Delete(c *Ctx) error {
	fileID := c.Path("fileID")

	file, err := h.repo.FileByID(c.Context(), fileID)
	if err != nil {
		return err
	}

	if !c.IsSpaceAdmin(c.User, file.Folder.Space) {
		return httperror.Forbidden
	}

	if err := h.repo.DeleteFile(c.Context(), fileID); err != nil {
		return err
	}

	// reload folder
	folder, err := h.repo.FolderByID(c.Context(), file.FolderID)
	if err != nil {
		return err
	}

	return c.HTML(
		http.StatusOK,
		"",
		"show_folder/refresh_show_files",
		Map{
			"folder": folder,
		},
	)
}
