// Copyright 2026 Lizandro
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package db

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/U-SITE-SAS-BIC/utasker/models"
)

const configDirName = ".tasker"
const tasksFileName = "tasks.json"
const projectFileName = ".task-project"

var mu sync.Mutex

func getConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("home dir: %w", err)
	}
	dir := filepath.Join(home, configDirName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("create config dir: %w", err)
	}
	return dir, nil
}

type storage struct {
	Tasks  []models.Task `json:"tasks"`
	NextID int           `json:"next_id"`
}

func loadTasks() (storage, error) {
	dir, err := getConfigDir()
	if err != nil {
		return storage{}, err
	}
	path := filepath.Join(dir, tasksFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return storage{NextID: 1}, nil
		}
		return storage{}, fmt.Errorf("read tasks: %w", err)
	}
	var s storage
	if err := json.Unmarshal(data, &s); err != nil {
		return storage{}, fmt.Errorf("parse tasks: %w", err)
	}
	return s, nil
}

func saveTasks(s storage) error {
	dir, err := getConfigDir()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, tasksFileName)
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal tasks: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write tasks: %w", err)
	}
	return nil
}

func GetProjectFromDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getwd: %w", err)
	}
	path := filepath.Join(cwd, projectFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("read project file: %w", err)
	}
	var cfg struct {
		Project string `json:"project"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return "", fmt.Errorf("parse project file: %w", err)
	}
	return cfg.Project, nil
}

func InitProject(project string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getwd: %w", err)
	}
	path := filepath.Join(cwd, projectFileName)
	cfg := struct {
		Project string `json:"project"`
	}{Project: project}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write project file: %w", err)
	}
	return nil
}

func AddTask(title, project, description string, priority int, tags []string, dueDate string) (models.Task, error) {
	mu.Lock()
	defer mu.Unlock()

	s, err := loadTasks()
	if err != nil {
		return models.Task{}, err
	}

	task := models.Task{
		ID:          s.NextID,
		Project:     project,
		Title:       title,
		Description: description,
		Priority:    priority,
		Status:      models.StatusPending,
		Tags:        tags,
		DueDate:     dueDate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	s.Tasks = append(s.Tasks, task)
	s.NextID++

	if err := saveTasks(s); err != nil {
		return models.Task{}, err
	}

	return task, nil
}

func ListTasks(project, status string, all bool) ([]models.Task, error) {
	mu.Lock()
	defer mu.Unlock()

	s, err := loadTasks()
	if err != nil {
		return nil, err
	}

	var result []models.Task
	for _, t := range s.Tasks {
		if !all && project != "" && t.Project != project {
			continue
		}
		if status != "" && t.Status != status {
			continue
		}
		result = append(result, t)
	}
	return result, nil
}

func GetTask(id int) (models.Task, error) {
	mu.Lock()
	defer mu.Unlock()

	s, err := loadTasks()
	if err != nil {
		return models.Task{}, err
	}

	for _, t := range s.Tasks {
		if t.ID == id {
			return t, nil
		}
	}
	return models.Task{}, fmt.Errorf("task %d not found", id)
}

func UpdateTaskStatus(id int, status string) (models.Task, error) {
	mu.Lock()
	defer mu.Unlock()

	s, err := loadTasks()
	if err != nil {
		return models.Task{}, err
	}

	for i, t := range s.Tasks {
		if t.ID == id {
			t.Status = status
			t.UpdatedAt = time.Now()
			s.Tasks[i] = t
			if err := saveTasks(s); err != nil {
				return models.Task{}, err
			}
			return t, nil
		}
	}
	return models.Task{}, fmt.Errorf("task %d not found", id)
}

func UpdateTask(id int, title, description string, priority int, tags []string, dueDate string) (models.Task, error) {
	mu.Lock()
	defer mu.Unlock()

	s, err := loadTasks()
	if err != nil {
		return models.Task{}, err
	}

	for i, t := range s.Tasks {
		if t.ID == id {
			if title != "" {
				t.Title = title
			}
			t.Description = description
			t.Priority = priority
			if tags != nil {
				t.Tags = tags
			}
			t.DueDate = dueDate
			t.UpdatedAt = time.Now()
			s.Tasks[i] = t
			if err := saveTasks(s); err != nil {
				return models.Task{}, err
			}
			return t, nil
		}
	}
	return models.Task{}, fmt.Errorf("task %d not found", id)
}

func DeleteTask(id int) error {
	mu.Lock()
	defer mu.Unlock()

	s, err := loadTasks()
	if err != nil {
		return err
	}

	for i, t := range s.Tasks {
		if t.ID == id {
			s.Tasks = append(s.Tasks[:i], s.Tasks[i+1:]...)
			return saveTasks(s)
		}
	}
	return fmt.Errorf("task %d not found", id)
}
