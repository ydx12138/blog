package dao

import (
	"blog/core"
	"blog/models"

	"go.uber.org/zap"
)

// GetAllCategories 获取所有分类（按sort排序）
func GetAllCategories() ([]models.Category, error) {
	var categories []models.Category = make([]models.Category, 0)
	err := core.DB.Order("sort DESC").Find(&categories).Error
	if err != nil {
		zap.L().Error("GetAllCategories:" + err.Error())
		return categories, err
	}
	return categories, nil
}

// GetCategoryByID 根据ID获取分类
func GetCategoryByID(id uint64) (models.Category, error) {
	var category models.Category
	err := core.DB.First(&category, id).Error
	if err != nil {
		zap.L().Error("GetCategoryByID:" + err.Error())
		return category, err
	}
	return category, nil
}

// CreateCategory 创建分类
func CreateCategory(category *models.Category) error {
	err := core.DB.Create(category).Error
	if err != nil {
		zap.L().Error("CreateCategory:" + err.Error())
		return err
	}
	return nil
}

// UpdateCategory 更新分类
func UpdateCategory(category *models.Category) error {
	err := core.DB.Save(category).Error
	if err != nil {
		zap.L().Error("UpdateCategory:" + err.Error())
		return err
	}
	return nil
}

// DeleteCategory 删除分类
func DeleteCategory(id uint64) error {
	err := core.DB.Delete(&models.Category{}, id).Error
	if err != nil {
		zap.L().Error("DeleteCategory:" + err.Error())
		return err
	}
	return nil
}
